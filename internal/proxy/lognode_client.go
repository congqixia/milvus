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

type LogNodeCreatorFunc func(ctx context.Context, addr string, nodeID int64) (types.LogNodeClient, error)

type LognodeManager interface {
	// get log node client by channel name
	GetChannelClient(channel string) (types.LogNodeClient, error)
	GetNodeClient(nodeID int64) (types.LogNodeClient, error)
	UpdateDistribution(ctx context.Context) error
	Close()
}

type nodeClient struct {
	nodeID   int64
	address  string
	isClosed bool
	client   types.LogNodeClient
}

func (c *nodeClient) close() error {
	if !c.isClosed {
		err := c.client.Close()
		if err != nil {
			return err
		}
		c.isClosed = true
	}
	return nil
}

func (c *nodeClient) getClient() (types.LogNodeClient, error) {
	if c.isClosed {
		return nil, errors.New("client is closed")
	}
	return c.client, nil
}

type LognodeManagerImpl struct {
	// nodeID -> NodeClient
	clients       map[UniqueID]*nodeClient
	clientCreator LogNodeCreatorFunc

	// pchannel -> NodeClient
	channelDistribution map[string]*nodeClient
	dc                  types.DataCoordClient

	mu sync.RWMutex
}

func (m *LognodeManagerImpl) getNodeClient(nodeID int64) (types.LogNodeClient, error) {
	client, ok := m.clients[nodeID]
	if !ok {
		return nil, fmt.Errorf("can not find client of node %d", nodeID)
	}
	return client.getClient()
}

func (m *LognodeManagerImpl) GetNodeClient(nodeID int64) (types.LogNodeClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.getNodeClient(nodeID)
}

func (m *LognodeManagerImpl) GetChannelClient(channel string) (types.LogNodeClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, ok := m.channelDistribution[channel]
	if !ok {
		return nil, merr.WrapErrChannelNotFound(channel, "channel log node not in distribution")
	}

	return client.getClient()
}

func (m *LognodeManagerImpl) updateClients(ctx context.Context, nodes []*logpb.LogNodeInfo) error {
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
			if err := client.close(); err != nil {
				log.Warn("close invalid log node client failed", zap.Int64("nodeID", client.nodeID), zap.Error(err))
				return err
			}
		}
	}
	m.clients = newClients
	return nil
}

func (m *LognodeManagerImpl) UpdateDistribution(ctx context.Context) error {
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
