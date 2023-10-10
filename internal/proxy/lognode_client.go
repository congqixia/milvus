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
	"github.com/milvus-io/milvus/internal/types"
	"github.com/milvus-io/milvus/pkg/util/merr"
)

type logNodeCreatorFunc func(ctx context.Context, addr string, nodeID int64) (types.LogNodeClient, error)

type lognodeManager interface {
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

func (c *nodeClient) getClient() (types.LogNodeClient, error) {
	if c.isClosed {
		return nil, errors.New("client is closed")
	}
	return c.client, nil
}

type lognodeManagerImpl struct {
	// nodeID -> nodeClient
	clients       map[UniqueID]*nodeClient
	clientCreator logNodeCreatorFunc

	// pchannel -> nodeClient
	channelDistribution map[string]*nodeClient
	rc                  types.RootCoordClient

	mu sync.RWMutex
}

func (m *lognodeManagerImpl) getNodeClient(nodeID int64) (types.LogNodeClient, error) {
	client, ok := m.clients[nodeID]
	if !ok {
		return nil, fmt.Errorf("can not find client of node %d", nodeID)
	}
	return client.getClient()
}

func (m *lognodeManagerImpl) GetNodeClient(nodeID int64) (types.LogNodeClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.getNodeClient(nodeID)
}

func (m *lognodeManagerImpl) GetChannelClient(channel string) (types.LogNodeClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, ok := m.channelDistribution[channel]
	if !ok {
		return nil, merr.WrapErrChannelNotFound(channel, "channel log node not in distribution")
	}

	return client.getClient()
}

func (m *lognodeManagerImpl) UpdateDistribution(ctx context.Context) error {
	return nil
}
