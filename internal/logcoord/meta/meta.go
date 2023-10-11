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

	"github.com/milvus-io/milvus/internal/metastore"
	"github.com/milvus-io/milvus/internal/proto/logpb"
	"github.com/milvus-io/milvus/internal/types"
	"github.com/samber/lo"
)

type Meta interface {
	Init(types.RootCoordClient) error
	// physical Channel
	GetPChannel(channel string) *PhysicalChannel
	GetPChannelInfos() []*logpb.PChannelInfo
	GetPChannelNamesBy(filters ...PChannelFilter) []string

	UpdateLeaseID(ctx context.Context, channel string) error
	AssignPChannel(ctx context.Context, channel string, nodeID int64) error
	UnassignPChannel(ctx context.Context, channel string) error

	// virtual vhannel
	AddVChannel(channels ...string) error
	RemoveVChannel(channels ...string) error
}

type PChannelFilter func(*PhysicalChannel) bool

func WithState(state logpb.PChannelState) PChannelFilter {
	return func(c *PhysicalChannel) bool {
		return c.CheckState(state)
	}
}

type PChannelList map[string]*PhysicalChannel

func (l *PChannelList) Get(channel string) *PhysicalChannel {
	return (*l)[channel]
}

func (l *PChannelList) GetNames() []string {
	return lo.Keys(*l)
}

func (l *PChannelList) GetInfos() []*logpb.PChannelInfo {
	return lo.Map(lo.Values(*l), func(channel *PhysicalChannel, _ int) *logpb.PChannelInfo {
		return channel.GetInfo()
	})
}

func NewPChannelList(channels []string, catalog metastore.DataCoordCatalog, infos map[string]*logpb.PChannelInfo, leaseIDs map[string]uint64) PChannelList {
	list := make(map[string]*PhysicalChannel)

	for _, channel := range channels {
		pChannel := NewPhysicalChannel(channel, catalog)
		info, ok := infos[channel]
		if ok {
			pChannel.nodeID = info.GetNodeID()
		}

		leaseID, ok := leaseIDs[channel]
		if ok {
			pChannel.leaseID = leaseID
		}
		list[channel] = pChannel
	}

	return list
}

type MetaImpl struct {
	catalog     metastore.DataCoordCatalog
	channelList PChannelList
}

func NewMetaImpl(catalog metastore.DataCoordCatalog) *MetaImpl {
	return &MetaImpl{
		catalog: catalog,
	}
}

func (m *MetaImpl) initPChannel(ctx context.Context, channels ...string) error {
	infos, err := m.catalog.ListPChannelInfo(ctx)
	if err != nil {
		return err
	}

	leaseIDs, err := m.catalog.ListPChannelLeaseID(ctx)
	if err != nil {
		return err
	}

	m.channelList = NewPChannelList(channels, m.catalog, infos, leaseIDs)
	return nil
}

func (m *MetaImpl) Init(rc types.RootCoordClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init pchannel info
	pChannels := getPChannelList()
	err := m.initPChannel(ctx, pChannels...)
	if err != nil {
		return err
	}
	// TODO RELOAD VCHANNEL NUM OF PCHNNEL
	return nil
}

func (m *MetaImpl) UpdateLeaseID(ctx context.Context, channel string) error {
	return m.GetPChannel(channel).UpdateLeaseID(ctx)
}

func (m *MetaImpl) AssignPChannel(ctx context.Context, channel string, nodeID int64) error {
	return m.GetPChannel(channel).Assign(ctx, nodeID)
}

func (m *MetaImpl) UnassignPChannel(ctx context.Context, channel string) error {
	return m.GetPChannel(channel).Unassign(ctx)
}

func (m *MetaImpl) GetPChannel(channel string) *PhysicalChannel {
	return m.channelList.Get(channel)
}

func (m *MetaImpl) GetPChannelInfos() []*logpb.PChannelInfo {
	return m.channelList.GetInfos()
}

func (m *MetaImpl) GetPChannelNamesBy(filters ...PChannelFilter) []string {
	channels := []string{}
	filter := func(channel *PhysicalChannel) bool {
		for _, filter := range filters {
			if !filter(channel) {
				return false
			}
		}
		return true
	}

	for name, channel := range m.channelList {
		if filter(channel) {
			channels = append(channels, name)
		}
	}
	return channels
}

func (m *MetaImpl) AddVChannel(channels ...string) error {
	for _, channel := range channels {
		pchannel := getPChannelName(channel)
		m.channelList.Get(pchannel).IncRef()
	}
	return nil
}

func (m *MetaImpl) RemoveVChannel(channels ...string) error {
	for _, channel := range channels {
		pchannel := getPChannelName(channel)
		m.channelList.Get(pchannel).DecRef()
	}
	return nil
}
