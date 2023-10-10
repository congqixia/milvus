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

package session

import (
	"context"
	"sync"

	"github.com/milvus-io/milvus/internal/logcoord/meta"
	"github.com/milvus-io/milvus/internal/proto/logpb"
	"github.com/milvus-io/milvus/pkg/log"
	"github.com/milvus-io/milvus/pkg/util/retry"
	"go.uber.org/zap"
)

type SessionBalancer struct {
	nodeAllocator  NodeAllocator
	sessionManager *SessionManager
	meta           meta.Meta

	balanceNotify chan struct{}

	wg        sync.WaitGroup
	startOnce sync.Once
	stopOnce  sync.Once
	stopCh    chan struct{}
}

func NewSessionBalancer(sessionManager *SessionManager, meta meta.Meta) *SessionBalancer {
	nodeAllocator := NewUniformNodeAllocator(meta)
	return &SessionBalancer{
		meta:           meta,
		nodeAllocator:  nodeAllocator,
		sessionManager: sessionManager,
		balanceNotify:  make(chan struct{}, 1),
		stopCh:         make(chan struct{}),
	}
}

func (ba *SessionBalancer) allocWaitting(ctx context.Context) {
	channels := ba.meta.GetPChannelNamesBy(meta.WithState(logpb.PChannelState_Waitting))
	log.Info("alloc lognode for waitting pchannels", zap.Strings("channels", channels))
	for _, channel := range channels {
		nodeID := ba.nodeAllocator.Alloc(channel)
		if nodeID == -1 {
			log.Warn("no avaliable log node")
			return
		}

		err := retry.Do(ctx, func() error {
			err := ba.sessionManager.AssignPChannel(ctx, channel, nodeID)
			return err
		}, retry.Attempts(3))

		if err != nil {
			ba.nodeAllocator.Unalloc(channel)
			ba.nodeAllocator.FreezeNode(nodeID)
		}
	}
}

func (ba *SessionBalancer) balance() {
	defer ba.wg.Done()

	log.Info("start log node balancer loop")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case <-ba.balanceNotify:
			ba.allocWaitting(ctx)
			//TODO REASSIGN CHANNEL TO NOT BALANCE NODE
		case <-ba.stopCh:
			log.Info("close log node balancer")
			return
		}
	}
}

func (ba *SessionBalancer) Start() {
	ba.startOnce.Do(func() {
		ba.wg.Add(1)
		go ba.balance()
		ba.Notify()
	})
}

func (ba *SessionBalancer) Stop() {
	ba.stopOnce.Do(func() {
		close(ba.stopCh)
		ba.wg.Wait()
	})
}

func (ba *SessionBalancer) AddNode(nodeID int64) {
	ba.nodeAllocator.AddNode(nodeID)
	ba.Notify()
}

func (ba *SessionBalancer) RemoveNode(nodeID int64) {
	offlineChannels := ba.nodeAllocator.RemoveNode(nodeID)
	for _, channel := range offlineChannels {
		// TODO CTX?
		ba.meta.UnassignPChannel(context.Background(), channel)
	}
	ba.Notify()
}

func (ba *SessionBalancer) Notify() {
	select {
	case ba.balanceNotify <- struct{}{}:
		log.Info("notify log node balancer to do balance")
	default:
	}
}
