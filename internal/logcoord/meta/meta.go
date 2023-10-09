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
)

type Meta interface {
	Init(types.RootCoordClient) error
	// physical Channel
	GetPChannel(channel string) *PhysicalChannel
	GetPChannelList() PChannelList
	UpdateLeaseID(ctx context.Context, channel string)
	AssignPChannel(ctx context.Context, channel string, nodeID int64) error
	UnassignPChannel(ctx context.Context, channel string) error

	// virtual vhannel
	AddVChannel(channels ...string) error
	RemoveVChannel(channels ...string) error
	// ListVChannelName() []string
}

type PChannelList map[string]*PhysicalChannel

func (l *PChannelList) Get(channel string) *PhysicalChannel {
	return (*l)[channel]
}

func NewPChannelList(channels []string, infos map[string]*logpb.PChannelInfo, leaseIDs map[string]uint64) PChannelList {
	list := make(map[string]*PhysicalChannel)

	for _, channel := range channels {
		pChannel := NewPhysicalChannel(channel)
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

type ChannelMeta struct {
	catalog     metastore.DataCoordCatalog
	channelList PChannelList
}

func NewChannelMeta(catalog metastore.DataCoordCatalog) *ChannelMeta {
	return &ChannelMeta{
		catalog: catalog,
	}
}

func (m *ChannelMeta) initPChannel(ctx context.Context, channels ...string) error {
	infos, err := m.catalog.ListPChannelInfo(ctx)
	if err != nil {
		return err
	}

	leaseIDs, err := m.catalog.ListPChannelLeaseID(ctx)
	if err != nil {
		return err
	}

	m.channelList = NewPChannelList(channels, infos, leaseIDs)
	return nil
}

func (m *ChannelMeta) Init(rc types.RootCoordClient) error {
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

func (m *ChannelMeta) UpdateLeaseID(ctx context.Context, channel string) error {
	return m.GetPChannel(channel).UpdateLeaseID(ctx)
}

func (m *ChannelMeta) AssignPChannel(ctx context.Context, channel string, nodeID int64) error {
	return m.GetPChannel(channel).Assign(ctx, nodeID)
}

func (m *ChannelMeta) UnassignPChannel(ctx context.Context, channel string, nodeID int64) error

func (m *ChannelMeta) GetPChannel(channel string) *PhysicalChannel {
	return m.channelList.Get(channel)
}

func (m *ChannelMeta) GetPChannelList() PChannelList {
	return m.channelList
}

func (m *ChannelMeta) AddVChannel(channels ...string) error {
	for _, channel := range channels {
		pchannel := getPChannelName(channel)
		m.channelList.Get(pchannel).IncRef()
	}
	return nil
}

func (m *ChannelMeta) RemoveVChannel(channels ...string) error {
	for _, channel := range channels {
		pchannel := getPChannelName(channel)
		m.channelList.Get(pchannel).DecRef()
	}
	return nil
}
