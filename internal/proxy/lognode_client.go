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

package proxy

import (
	"context"
	"fmt"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/milvus-io/milvus/internal/proto/datapb"
	"github.com/milvus-io/milvus/internal/proto/logpb"
	"github.com/milvus-io/milvus/internal/types"
	"github.com/milvus-io/milvus/pkg/log"
	"github.com/milvus-io/milvus/pkg/util/merr"
	"go.uber.org/zap"
)

type logNodeCreatorFunc func(ctx context.Context, addr string, nodeID int64) (types.LogNodeClient, error)

type lognodeManager interface {
	// get log node client by channel name
	GetChannelClient(channel string) (int64, types.LogNodeClient, error)
	GetNodeClient(nodeID int64) (types.LogNodeClient, error)
	UpdateDistribution(ctx context.Context) error
	Close()
}

type nodeClient struct {
	nodeID  int64
	address string
	client  types.LogNodeClient
}

func (c *nodeClient) Close() error {
	if c.client != nil {
		err := c.client.Close()
		if err != nil {
			return err
		}
		c.client = nil
	}
	return nil
}

func (c *nodeClient) GetNodeID() int64 {
	return c.nodeID
}

func (c *nodeClient) GetClient() (types.LogNodeClient, error) {
	if c.client == nil {
		return nil, errors.New("client not exist, may closed")
	}
	return c.client, nil
}

type lognodeManagerImpl struct {
	// nodeID -> nodeClient
	clients       map[UniqueID]*nodeClient
	clientCreator logNodeCreatorFunc

	// pchannel -> nodeClient
	channelDistribution map[string]*nodeClient
	dc                  types.DataCoordClient

	mu sync.RWMutex
}

func (m *lognodeManagerImpl) GetNodeClient(nodeID int64) (types.LogNodeClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, ok := m.clients[nodeID]
	if !ok {
		return nil, fmt.Errorf("can not find client of node %d", nodeID)
	}
	return client.GetClient()
}

func (m *lognodeManagerImpl) GetChannelClient(channel string) (int64, types.LogNodeClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	node, ok := m.channelDistribution[channel]
	if !ok {
		return 0, nil, merr.WrapErrChannelNotFound(channel, "channel log node not in distribution")
	}

	client, err := node.GetClient()
	if err != nil {
		return 0, nil, merr.WrapErrNodeNotAvailable(node.GetNodeID(), "get node client failed")
	}

	return node.GetNodeID(), client, nil
}

func (m *lognodeManagerImpl) updateClients(ctx context.Context, nodes []*logpb.LogNodeInfo) error {
	newClients := make(map[UniqueID]*nodeClient)
	for _, node := range nodes {
		if client, ok := m.clients[node.NodeID]; ok {
			newClients[node.GetNodeID()] = client
		} else {
			// TODO RETRY?
			client, err := m.clientCreator(ctx, node.GetAddress(), node.GetNodeID())
			if err != nil {
				log.Warn("connect new log node failed", zap.Int64("nodeID", node.GetNodeID()), zap.Error(err))
				return err
			}
			newClients[node.GetNodeID()] = &nodeClient{
				nodeID:  node.GetNodeID(),
				address: node.GetAddress(),
				client:  client,
			}
		}
	}

	for _, client := range m.clients {
		if client, ok := newClients[client.nodeID]; !ok {
			if err := client.Close(); err != nil {
				log.Warn("close invalid log node client failed", zap.Int64("nodeID", client.nodeID), zap.Error(err))
				return err
			}
		}
	}
	m.clients = newClients
	return nil
}

func (m *lognodeManagerImpl) UpdateDistribution(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	resp, err := m.dc.GetChannelDistribution(ctx, &datapb.GetChannelDistributionRequest{})
	if err != nil {
		return err
	}

	err = merr.Error(resp.GetStatus())
	if err != nil {
		return err
	}

	err = m.updateClients(ctx, resp.NodeInfos)
	if err != nil {
		return err
	}

	m.channelDistribution = make(map[string]*nodeClient)
	for _, channel := range resp.ChannelInfos {
		if client, ok := m.clients[channel.GetNodeID()]; ok {
			m.channelDistribution[channel.GetName()] = client
		} else {
			log.Warn("channel waitting assign or assign to unknown node",
				zap.String("channel", channel.GetName()),
				zap.Int64("nodeID", channel.GetNodeID()))
		}
	}

	return nil
}
