// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/blang/semver/v4"
	"github.com/cockroachdb/errors"
	"github.com/milvus-io/milvus/internal/registry/common"
	"github.com/milvus-io/milvus/internal/registry/options"
	milvuscommon "github.com/milvus-io/milvus/pkg/common"
	"github.com/milvus-io/milvus/pkg/log"
	"github.com/milvus-io/milvus/pkg/util/merr"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/atomic"
)

// var _ registry.Session = (*etcdSession)(nil)

type etcdSession struct {
	common.ServiceEntryBase

	idOnce sync.Once

	cancel context.CancelFunc

	// session info
	exclusive       bool
	stopping        *atomic.Bool
	version         semver.Version
	client          *clientv3.Client
	leaseID         clientv3.LeaseID
	isStandby       *atomic.Bool
	keepaliveCancel context.CancelFunc
	triggerKill     bool
	liveCh          <-chan struct{}
	closeCh         chan struct{}

	// options
	metaRoot            string
	useCustomConfig     bool
	sessionTTL          int64
	sessionRetryTimes   int64
	reuseNodeID         bool
	enableActiveStandBy bool
}

func newEtcdSession(cli *clientv3.Client, metaRoot string, addr string, component string,
	opt options.RegisterOpt) *etcdSession {
	return &etcdSession{
		ServiceEntryBase: common.NewServiceEntryBase(addr, component),
		exclusive:        opt.Exclusive,
		client:           cli,
		metaRoot:         metaRoot,
		isStandby:        atomic.NewBool(opt.StandBy),
		version:          milvuscommon.Version,
		stopping:         atomic.NewBool(false),
	}
}

func newEtcdEntry(addr string, service string, opt options.RegisterOpt) *etcdSession {
	base := common.NewServiceEntryBase(addr, service)
	return &etcdSession{
		ServiceEntryBase: base,
		exclusive:        opt.Exclusive,
		isStandby:        atomic.NewBool(opt.StandBy),
		version:          milvuscommon.Version,
		stopping:         atomic.NewBool(false),
		closeCh:          make(chan struct{}),
	}
}

// UnmarshalJSON unmarshal bytes to Session.
func (s *etcdSession) UnmarshalJSON(data []byte) error {
	var raw struct {
		ServerID   int64  `json:"ServerID,omitempty"`
		ServerName string `json:"ServerName,omitempty"`
		Address    string `json:"Address,omitempty"`
		Exclusive  bool   `json:"Exclusive,omitempty"`
		Stopping   bool   `json:"Stopping,omitempty"`
		Version    string `json:"Version"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	if raw.Version != "" {
		s.version, err = semver.Parse(raw.Version)
		if err != nil {
			return err
		}
	}

	s.SetID(raw.ServerID)
	s.SetAddr(raw.Address)
	s.SetComponentType(raw.ServerName)
	s.exclusive = raw.Exclusive
	s.stopping = atomic.NewBool(raw.Stopping)
	return nil
}

// MarshalJSON marshals session to bytes.
func (s *etcdSession) MarshalJSON() ([]byte, error) {

	verStr := s.version.String()
	return json.Marshal(&struct {
		ServerID    int64  `json:"ServerID,omitempty"`
		ServerName  string `json:"ServerName,omitempty"`
		Address     string `json:"Address,omitempty"`
		Exclusive   bool   `json:"Exclusive,omitempty"`
		Stopping    bool   `json:"Stopping,omitempty"`
		TriggerKill bool
		Version     string `json:"Version"`
	}{
		ServerID:    s.ID(),
		ServerName:  s.ComponentType(),
		Address:     s.Addr(),
		Exclusive:   s.exclusive,
		Stopping:    s.stopping.Load(),
		TriggerKill: s.triggerKill,
		Version:     verStr,
	})

}

func (s *etcdSession) Revoke(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if s.client == nil || s.leaseID == 0 {
		// TODO audit error type here
		return merr.WrapErrParameterInvalid("valid session", "session not valid")
	}

	_, err := s.client.Revoke(ctx, s.leaseID)
	if err != nil {
		if errors.IsAny(err, context.Canceled, context.DeadlineExceeded) {
			return errors.Wrapf(err, "context canceled when revoking component %s serverID:%d", s.ComponentType(), s.ID())
		}
		return merr.WrapErrIoFailed(fmt.Sprintf("%s-%d", s.ComponentType(), s.ID()), err.Error())
	}

	return nil
}

// Stop marks session as `Stopping` state for graceful stop.
/*
func (s *etcdSession) Stop(ctx context.Context) error {
	if s == nil || s.client == nil || s.leaseID == 0 {
		return merr.WrapErrParameterInvalid("valid session", "session not valid")
	}

	completeKey := s.getCompleteKey()
	log := log.Ctx(ctx).With(
		zap.String("key", completeKey),
	)
	resp, err := s.client.Get(ctx, completeKey, clientv3.WithCountOnly())
	if err != nil {
		log.Error("fail to get the session", zap.Error(err))
		return err
	}
	if resp.Count == 0 {
		return nil
	}
	s.stopping.Store(true)
	sessionJSON, err := json.Marshal(s)
	if err != nil {
		log.Error("fail to marshal the session")
		return err
	}
	_, err = s.client.Put(ctx, completeKey, string(sessionJSON), clientv3.WithLease(s.leaseID))
	if err != nil {
		log.Error("fail to update the session to stopping state")
		return err
	}
	return nil
}*/

func (e *etcdSession) key() string {
	return fmt.Sprintf("%s_%d", e.ComponentType(), e.ID())
}

// keepAlive processes the response of etcd keepAlive interface
// If keepAlive fails for unexpected error, it will send a signal to the channel.
func (s *etcdSession) keepAlive(ch <-chan *clientv3.LeaseKeepAliveResponse) (failChannel <-chan struct{}) {
	failCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-s.closeCh:
				log.Info("etcd session quit")
				return
			case resp, ok := <-ch:
				if !ok {
					log.Warn("session keepalive channel closed")
					close(failCh)
					return
				}
				if resp == nil {
					log.Warn("session keepalive response failed")
					close(failCh)
					return
				}
			}
		}
	}()
	return failCh
}
