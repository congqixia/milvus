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

package meta

import (
	"context"
	"sync"

	"github.com/milvus-io/milvus/internal/metastore"
	"github.com/milvus-io/milvus/internal/proto/logpb"
	"github.com/milvus-io/milvus/pkg/log"
	"go.uber.org/zap"
)

type PhysicalChannel struct {
	catalog metastore.DataCoordCatalog
	name    string
	state   logpb.PChannelState

	nodeID  int64
	leaseID uint64
	// count of vchannel
	refCnt uint64

	mu sync.RWMutex
}

func NewPhysicalChannel(name string, catalog metastore.DataCoordCatalog) *PhysicalChannel {
	return &PhysicalChannel{
		name:    name,
		nodeID:  -1,
		leaseID: 0,
		catalog: catalog,
		state:   logpb.PChannelState_Waitting,
	}
}

func (c *PhysicalChannel) Assign(ctx context.Context, nodeID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.catalog.SavePChannelInfo(ctx, &logpb.PChannelInfo{
		Name:   c.name,
		NodeID: nodeID,
		State:  logpb.PChannelState_Watching,
	})
	if err != nil {
		log.Warn("assign pchannel update etcd channel info failed", zap.String("channel", c.name), zap.Error(err))
		return err
	}

	c.nodeID = nodeID
	c.state = logpb.PChannelState_Watching
	return nil
}

func (c *PhysicalChannel) Unassign(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.catalog.SavePChannelInfo(ctx, &logpb.PChannelInfo{
		Name:   c.name,
		NodeID: -1,
		State:  logpb.PChannelState_Waitting,
	})
	if err != nil {
		log.Warn("assign pchannel update etcd channel info failed", zap.String("channel", c.name), zap.Error(err))
		return err
	}

	c.nodeID = -1
	c.state = logpb.PChannelState_Waitting
	return nil
}

func (c *PhysicalChannel) UpdateLeaseID(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.catalog.SavePChannelLeaseID(ctx, c.name, c.leaseID+1)
	if err != nil {
		log.Warn("pchannel update etcd channel leaseID failed", zap.String("channel", c.name), zap.Error(err))
		return err
	}
	c.leaseID = c.leaseID + 1
	return nil
}

func (c *PhysicalChannel) GetNodeID() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.nodeID
}

func (c *PhysicalChannel) GetLeaseID() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.leaseID
}

func (c *PhysicalChannel) GetRef() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.refCnt
}

func (c *PhysicalChannel) GetInfo() *logpb.PChannelInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return &logpb.PChannelInfo{
		Name:   c.name,
		State:  c.state,
		NodeID: c.nodeID,
	}
}

func (c *PhysicalChannel) CheckState(state logpb.PChannelState) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.state == state
}

func (c *PhysicalChannel) IncRef() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.refCnt++
}

func (c *PhysicalChannel) DecRef() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.refCnt--
}

func (c *PhysicalChannel) IsUsed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.refCnt > 0
}
