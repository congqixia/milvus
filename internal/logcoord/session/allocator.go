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
	"sync"
	"time"

	"github.com/milvus-io/milvus/internal/logcoord/meta"
	"github.com/milvus-io/milvus/pkg/log"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// NodeAllocator alloc log node for pchannel
type NodeAllocator interface {
	Alloc(pChannel string) int64
	Unalloc(pChannel string)
	//create a realloc plan for balance
	//empty string means don't need balance
	//return pchannel, node, target_nodes
	Realloc() (string, int64, int64)

	AddNode(nodeID int64)
	RemoveNode(nodeID int64) []string
	FreezeNode(nodeID int64)
}

const FreezeInterval = 10 * time.Second

type NodeInfo struct {
	// pChannel list
	channelList map[string]*meta.PhysicalChannel
	FreezeTime  time.Time
	meta        meta.Meta
}

func NewNodeInfo() *NodeInfo {
	return &NodeInfo{
		channelList: make(map[string]*meta.PhysicalChannel),
	}
}

func (i *NodeInfo) RemoveChannel(channel string) {
	delete(i.channelList, channel)
}

func (i *NodeInfo) GetChannelNum() int64 {
	return int64(len(i.channelList))
}

func (i *NodeInfo) GetChannelNames() []string {
	return lo.Keys(i.channelList)
}

func (i *NodeInfo) AddChannel(channel string) {
	if _, ok := i.channelList[channel]; !ok {
		i.channelList[channel] = i.meta.GetPChannel(channel)
	}
}

func (i *NodeInfo) PopChannel() string {
	for name, _ := range i.channelList {
		delete(i.channelList, name)
		return name
	}

	return ""
}

func (i *NodeInfo) IsFrozen() bool {
	return time.Now().Before(i.FreezeTime)
}

type UniformNodeAllocator struct {
	// nodeID -> NodeInfo
	nodeInfos map[int64]*NodeInfo
	// pChannel -> nodeID
	nodeMapping map[string]int64

	mu sync.Mutex
}

func (allocator *UniformNodeAllocator) AddNode(nodeID int64) {
	allocator.mu.Lock()
	defer allocator.mu.Unlock()

	_, ok := allocator.nodeInfos[nodeID]
	if !ok {
		allocator.nodeInfos[nodeID] = NewNodeInfo()
	}
}

func (allocator *UniformNodeAllocator) RemoveNode(nodeID int64) []string {
	allocator.mu.Lock()
	defer allocator.mu.Unlock()
	info, ok := allocator.nodeInfos[nodeID]
	if !ok {
		return []string{}
	}

	if info.GetChannelNum() != 0 {
		log.Info("Remove log node but not reassign all pchannel", zap.Int64("nodeID", nodeID))
	}
	delete(allocator.nodeInfos, nodeID)
	return info.GetChannelNames()
}

func (allocator *UniformNodeAllocator) FreezeNode(nodeID int64) {
	allocator.mu.Lock()
	defer allocator.mu.Unlock()

	info, ok := allocator.nodeInfos[nodeID]
	if !ok {
		return
	}

	info.FreezeTime = time.Now().Add(FreezeInterval)
}

func (allocator *UniformNodeAllocator) selectMinNode() (int64, int64) {
	minNode, minCount := int64(-1), int64(-1)
	for nodeID, info := range allocator.nodeInfos {
		if info.IsFrozen() {
			continue
		}

		count := info.GetChannelNum()
		if minNode == -1 || count < minCount {
			minNode = nodeID
			minCount = count
		}
	}
	return minNode, minCount
}

func (allocator *UniformNodeAllocator) selectMaxNode() (int64, int64) {
	maxNode, maxCount := int64(-1), int64(-1)
	for nodeID, info := range allocator.nodeInfos {
		if info.IsFrozen() {
			continue
		}

		count := info.GetChannelNum()
		if maxNode == -1 || count > maxCount {
			maxNode = nodeID
			maxCount = count
		}
	}
	return maxNode, maxCount
}

func (allocator *UniformNodeAllocator) Alloc(pChannel string) int64 {
	allocator.mu.Lock()
	defer allocator.mu.Unlock()

	minNode, _ := allocator.selectMinNode()
	allocator.nodeInfos[minNode].AddChannel(pChannel)
	allocator.nodeMapping[pChannel] = minNode
	return minNode
}

func (allocator *UniformNodeAllocator) Unalloc(pChannel string) {
	allocator.mu.Lock()
	defer allocator.mu.Unlock()

	nodeID, ok := allocator.nodeMapping[pChannel]
	if ok {
		allocator.nodeInfos[nodeID].RemoveChannel(pChannel)
		delete(allocator.nodeMapping, pChannel)
	}
}

func (allocator *UniformNodeAllocator) Realloc() (string, int64, int64) {
	allocator.mu.Lock()
	defer allocator.mu.Unlock()

	maxNode, maxCount := allocator.selectMaxNode()
	minNode, minCount := allocator.selectMinNode()
	if maxCount <= minCount+2 || maxCount == 0 {
		return "", -1, -1
	}

	target_channel := allocator.nodeInfos[maxNode].PopChannel()
	allocator.nodeMapping[target_channel] = minNode
	return target_channel, maxNode, minNode
}

func NewUniformNodeAllocator() *UniformNodeAllocator {
	return &UniformNodeAllocator{
		nodeInfos:   make(map[int64]*NodeInfo),
		nodeMapping: make(map[string]int64),
	}
}
