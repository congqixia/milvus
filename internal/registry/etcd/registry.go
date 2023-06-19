// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the Licensr. You may obtain a copy of the License at
//
//     http://www.apachr.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the Licensr.

package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strconv"

	"github.com/milvus-io/milvus/internal/registry"
	"github.com/milvus-io/milvus/internal/registry/options"
	"github.com/milvus-io/milvus/pkg/common"
	"github.com/milvus-io/milvus/pkg/log"
	"github.com/milvus-io/milvus/pkg/util/merr"
	"github.com/milvus-io/milvus/pkg/util/retry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// etcdRegistry implements registry.Registry with Etcd.
type etcdRegistry struct {
	client   *clientv3.Client
	metaRoot string
	manager  *serviceManager
}

func NewEtcdRegistry(ctx context.Context, cli *clientv3.Client, metaRoot string) (registry.Registry, error) {
	r := &etcdRegistry{
		client:   cli,
		metaRoot: metaRoot,
		manager:  newServiceManager(cli),
	}
	err := r.initIDCounter(ctx)
	if err != nil {
		log.Warn("failed to initialize serverID counter", zap.Error(err))
		return nil, err
	}

	return r, nil
}

func (r *etcdRegistry) RegisterService(ctx context.Context, service string, addr string, opts ...options.RegisterOption) (registry.ServiceEntry, error) {

	opt := options.DefaultRegisterOpt()
	for _, o := range opts {
		o(&opt)
	}

	session := newEtcdEntry(addr, service, opt)

	err := r.registerService(ctx, service, session, opt)
	if err != nil {
		return nil, err
	}
	return session, err
}

func (r *etcdRegistry) registerService(ctx context.Context, service string, session *etcdSession, opt options.RegisterOpt) error {
	log := log.Ctx(ctx).With(
		zap.String("service", service),
		zap.Bool("preset", false),
	)

	// TODO nodeID pre-set
	serverID, err := r.getServerID(ctx, service)
	if err != nil {
		return err
	}

	session.SetID(serverID)
	log.Info("get server id", zap.Int64("serverID", serverID))

	var key string
	if !opt.Exclusive {
		key = fmt.Sprintf("%s-%d", service, serverID) // non-exclusive, key shall be {serviceName}-{serverID}
	} else {
		key = service // exclusive, use only {serviceName} key
	}

	return r.grantSession(ctx, key, session, opt)
}

// grantSession performs grant operation for session kv in etcd.
func (r *etcdRegistry) grantSession(ctx context.Context, key string, session *etcdSession, opt options.RegisterOpt) error {
	key = path.Join(r.metaRoot, common.DefaultServiceRoot, key)

	fn := func() error {
		resp, err := r.client.Grant(ctx, opt.SessionTTL)
		if err != nil {
			log.Warn("grant session failed", zap.Error(err))
			return err
		}

		session.leaseID = resp.ID
		bs, err := json.Marshal(session)
		if err != nil {
			return retry.Unrecoverable(err)
		}

		txnResp, err := r.client.Txn(ctx).If(
			clientv3.Compare(clientv3.Version(key), "=", 0)).
			Then(clientv3.OpPut(key, string(bs))).Commit()
		if err != nil {
			return err
		}

		// TODO change to comparable error
		if !txnResp.Succeeded {
			return errors.New("cas failed")
		}

		return r.manager.Add(ctx, session)
	}

	err := retry.Do(ctx, fn, retry.Attempts(10))
	if err != nil {
		return err
	}

	return nil
}

func (r *etcdRegistry) getServerID(ctx context.Context, service string) (int64, error) {
	counterPath := path.Join(r.metaRoot, common.DefaultServiceRoot, common.DefaultIDKey)
	log := log.Ctx(ctx).With(
		zap.String("keyPath", counterPath),
		zap.String("service", service),
	)
	for {
		resp, err := r.client.Get(ctx, counterPath)
		if err != nil {
			log.Warn("etcdRegistry failed to read idCounter", zap.Error(err))
			return -1, merr.WrapErrIoFailed(counterPath, "etcd registry read id counter failed", err.Error())
		}
		if resp.Count == 0 {
			log.Warn("etcdRegistry failed to find idCounter")
			// should not retry, initialization shall be done before
			return -1, merr.WrapErrIoKeyNotFound(counterPath, "idCounter key not found")
		}

		value := string(resp.Kvs[0].Value)
		counter, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Warn("etcdRegistry failed to parse idCounter", zap.String("value", value))
			return -1, merr.WrapErrParameterInvalid("number id counter", value)
		}

		txnResp, err := r.client.Txn(ctx).If(
			clientv3.Compare(
				clientv3.Value(counterPath),
				"=",
				value)).
			Then(clientv3.OpPut(counterPath, strconv.FormatInt(counter+1, 10))).Commit()
		if err != nil {
			log.Warn("etcdRegistry txn add counter faild", zap.Error(err))
			return -1, merr.WrapErrIoFailed(counterPath, "failed to incr idCounter", err.Error())
		}

		if !txnResp.Succeeded {
			log.Warn("etcdRegistry get serverID raced, backoff")
			continue
		}

		return counter, nil
	}
}

// initIDCounter initializes server id counter.
func (r *etcdRegistry) initIDCounter(ctx context.Context) error {
	counterPath := path.Join(r.metaRoot, common.DefaultServiceRoot, common.DefaultIDKey)
	resp, err := r.client.Txn(ctx).If(
		clientv3.Compare(
			clientv3.Version(counterPath),
			"=",
			0)).
		Then(clientv3.OpPut(counterPath, "1")).Commit()
	if err != nil {
		return merr.WrapErrIoFailed(counterPath, "failed to initialize idCounter", err.Error())
	}
	// compare succeed, initialization completed
	if resp.Succeeded {
		log.Info("server ID counter initialized")
	}
	return nil
}

func (r *etcdRegistry) DeregisterService(ctx context.Context, entry registry.ServiceEntry) error {
	return r.manager.Remove(ctx, entry)
}

func (r *etcdRegistry) GetService(ctx context.Context, service string) (registry.ServiceEntry, error) {
	key := path.Join(r.metaRoot, common.DefaultServiceRoot, service)
	log.Ctx(ctx).With(
		zap.String("key", key),
	)
	resp, err := r.client.Get(ctx, key)
	if err != nil {
		log.Warn("failed to get service session", zap.Error(err))
		return nil, merr.WrapErrIoFailed(key, "failed to get service session", err.Error())
	}

	if resp.Count == 0 {
		log.Warn("no service online")
		return nil, merr.WrapErrIoKeyNotFound(key, "no service online")
	}

	session := &etcdSession{}
	err = json.Unmarshal(resp.Kvs[0].Value, session)
	if err != nil {
		log.Warn("failed to unmarhal service", zap.Error(err))
		return nil, merr.WrapErrParameterInvalid("valid session", string(resp.Kvs[0].Value))
	}

	// TODO setup options
	return newEtcdEntry(session.Addr(), service, options.DefaultRegisterOpt()), nil
}

func (r *etcdRegistry) ListServices(ctx context.Context, service string) ([]registry.ServiceEntry, error) {
	prefix := path.Join(r.metaRoot, common.DefaultServiceRoot, service)
	log.Ctx(ctx).With(
		zap.String("service", service),
		zap.String("key", prefix),
	)
	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Warn("failed to get service session", zap.Error(err))
		return nil, merr.WrapErrIoFailed(prefix, "failed to get service session", err.Error())
	}

	result := make([]registry.ServiceEntry, 0, resp.Count)
	for _, kv := range resp.Kvs {
		session := &etcdSession{}
		err = json.Unmarshal(kv.Value, session)
		if err != nil {
			log.Warn("failed to unmarhal service", zap.Error(err))
			return nil, merr.WrapErrParameterInvalid("valid session", string(resp.Kvs[0].Value))
		}
		// TODO setup options
		result = append(result, newEtcdEntry(session.Addr(), service, options.DefaultRegisterOpt()))
	}

	return result, nil
}

func (r *etcdRegistry) WatchService(ctx context.Context, service string, opts ...options.WatchOption) (registry.ServiceWatcher[any], error) {
	panic("not implemented") // TODO: Implement
}
