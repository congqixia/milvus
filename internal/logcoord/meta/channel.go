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
	"sync"

	"github.com/milvus-io/milvus/internal/proto/logpb"
)

type PhysicalChannel struct {
	name   string
	status logpb.PChannelState

	nodeID  int64
	leaseID uint64
	// count of vchannel
	refCnt uint64

	mu sync.RWMutex
}

func NewPhysicalChannel(name string) *PhysicalChannel {
	return &PhysicalChannel{
		name:    name,
		nodeID:  -1,
		leaseID: 0,
		status:  logpb.PChannelState_Waitting,
	}
}

func (c *PhysicalChannel) SetNodeID(nodeID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.nodeID = nodeID
}

func (c *PhysicalChannel) GetNodeID() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.nodeID
}

func (c *PhysicalChannel) SetLeaseID(leaseID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.leaseID = leaseID
}

func (c *PhysicalChannel) GetLeaseID() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.leaseID
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

// type VirtualChannel struct {
// 	name            string
// 	collectionID    uint64
// 	createTimestamp uint64
// 	startPosition   *msgpb.MsgPosition
// }

// func NewVirtualChannel(collectionID, createTimestamp uint64, name string, startPosition *msgpb.MsgPosition) *VirtualChannel {
// 	return &VirtualChannel{
// 		name:            name,
// 		collectionID:    collectionID,
// 		createTimestamp: createTimestamp,
// 		startPosition:   startPosition,
// 	}
// }
