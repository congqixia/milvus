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
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/milvus-io/milvus-proto/go-api/v2/commonpb"
	"github.com/milvus-io/milvus-proto/go-api/v2/milvuspb"
	"github.com/milvus-io/milvus-proto/go-api/v2/msgpb"
	"github.com/milvus-io/milvus/internal/proto/logpb"
	"github.com/milvus-io/milvus/internal/types"
	"github.com/milvus-io/milvus/pkg/log"
	"github.com/milvus-io/milvus/pkg/metrics"
	"github.com/milvus-io/milvus/pkg/mq/msgstream"
	"github.com/milvus-io/milvus/pkg/util/commonpbutil"
	"github.com/milvus-io/milvus/pkg/util/merr"
	"github.com/milvus-io/milvus/pkg/util/paramtable"
	"github.com/milvus-io/milvus/pkg/util/retry"
)

// channelsMgr manages the pchans, vchans
// repack and produce messages
type channelsMgr interface {
	GetChannels(collectionID UniqueID) ([]pChan, error)
	GetVChannels(collectionID UniqueID) ([]vChan, error)

	RemoveChannels(collectionID UniqueID)
	RemoveAllChannel()

	UpdateNodeInfo(ctx context.Context) error
	Insert(ctx context.Context, msgs map[string][]*msgpb.InsertRequest) error
}

type channelInfos struct {
	// It seems that there is no need to maintain relationships between vchans & pchans.
	vchans []vChan
	pchans []pChan
}

func newChannels(vchans []vChan, pchans []pChan) (channelInfos, error) {
	if len(vchans) != len(pchans) {
		log.Error("physical channels mismatch virtual channels", zap.Int("len(VirtualChannelNames)", len(vchans)), zap.Int("len(PhysicalChannelNames)", len(pchans)))
		return channelInfos{}, fmt.Errorf("physical channels mismatch virtual channels, len(VirtualChannelNames): %v, len(PhysicalChannelNames): %v", len(vchans), len(pchans))
	}

	return channelInfos{vchans: vchans, pchans: pchans}, nil
}

// getChannelsFuncType returns the channel information according to the collection id.
type getChannelsFuncType = func(collectionID UniqueID) (channelInfos, error)

// repackFuncType repacks message into message pack.
type repackFuncType = func(tsMsgs []msgstream.TsMsg, hashKeys [][]int32) (map[int32]*msgstream.MsgPack, error)

// getChannelsFunc returns a function about how to get channels of a collection.
func getChannelsFunc(ctx context.Context, rc types.RootCoordClient) getChannelsFuncType {
	return func(collectionID UniqueID) (channelInfos, error) {
		req := &milvuspb.DescribeCollectionRequest{
			Base:         commonpbutil.NewMsgBase(commonpbutil.WithMsgType(commonpb.MsgType_DescribeCollection)),
			CollectionID: collectionID,
		}

		resp, err := rc.DescribeCollection(ctx, req)
		if err != nil {
			log.Error("failed to describe collection", zap.Error(err), zap.Int64("collection", collectionID))
			return channelInfos{}, err
		}

		if resp.GetStatus().GetErrorCode() != commonpb.ErrorCode_Success {
			log.Error("failed to describe collection",
				zap.String("error_code", resp.GetStatus().GetErrorCode().String()),
				zap.String("reason", resp.GetStatus().GetReason()))
			return channelInfos{}, merr.Error(resp.GetStatus())
		}

		return newChannels(resp.GetVirtualChannelNames(), resp.GetPhysicalChannelNames())
	}
}

type channelsMgrImpl struct {
	nodes lognodeManager

	mu    sync.RWMutex              // mutex of channel infos
	infos map[UniqueID]channelInfos // collection id -> channel infos

	getChannelsFunc getChannelsFuncType
	repackFunc      repackFuncType
}

// NewChannelsMgrImpl constructs a channels manager.
func newChannelsMgrImpl(
	getDmlChannelsFunc getChannelsFuncType,
	dmlRepackFunc repackFuncType,
	msgStreamFactory msgstream.Factory,
) *channelsMgrImpl {
	return &channelsMgrImpl{
		infos:           make(map[UniqueID]channelInfos),
		getChannelsFunc: getDmlChannelsFunc,
		repackFunc:      dmlRepackFunc,
	}
}

func (mgr *channelsMgrImpl) GetAllChannels(collectionID UniqueID) (channelInfos, error) {
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()

	infos, ok := mgr.infos[collectionID]
	if ok {
		return infos, nil
	}

	info, err := mgr.getChannelsFunc(collectionID)
	if err != nil {
		return channelInfos{}, err
	}

	mgr.infos[collectionID] = info
	return info, nil
}

// getChannels returns the physical channels.
func (mgr *channelsMgrImpl) GetChannels(collectionID UniqueID) ([]pChan, error) {
	var channelInfos channelInfos
	channelInfos, err := mgr.GetAllChannels(collectionID)
	if err != nil {
		return nil, err
	}
	return channelInfos.pchans, nil
}

// getVChannels returns the virtual channels.
func (mgr *channelsMgrImpl) GetVChannels(collectionID UniqueID) ([]vChan, error) {
	var channelInfos channelInfos
	channelInfos, err := mgr.GetAllChannels(collectionID)
	if err != nil {
		return nil, err
	}
	return channelInfos.vchans, nil
}

// removeChannels remove channel of the specified collection. Idempotent.
// If channel already exists, remove it, otherwise do nothing.
func (mgr *channelsMgrImpl) RemoveChannels(collectionID UniqueID) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if info, ok := mgr.infos[collectionID]; ok {
		decPChanMetrics(info.pchans)
		delete(mgr.infos, collectionID)
	}
	log.Info("dml stream removed", zap.Int64("collection_id", collectionID))
}

// removeAllChannel remove all channel.
func (mgr *channelsMgrImpl) RemoveAllChannel() {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	for _, info := range mgr.infos {
		decPChanMetrics(info.pchans)
	}
	mgr.infos = make(map[UniqueID]channelInfos)
	log.Info("all dml stream removed")
}

func (mgr *channelsMgrImpl) UpdateNodeInfo(ctx context.Context) error {
	return mgr.nodes.UpdateDistribution(ctx)
}

func (mgr *channelsMgrImpl) insert(ctx context.Context, channel string, msgs []*msgpb.InsertRequest) error {
	nodeID, client, err := mgr.nodes.GetChannelClient(channel)
	if err != nil {
		log.Warn("Insert msg to channel failed, channel client not found", zap.String("channel", channel))
		return err
	}

	resp, err := client.Insert(ctx, &logpb.InsertRequest{
		Base: commonpbutil.NewMsgBase(
			commonpbutil.WithSourceID(paramtable.GetNodeID()),
			commonpbutil.WithTargetID(nodeID),
		),
		Msgs:      msgs,
		PChannels: []string{channel},
	})

	if err != nil {
		return err
	}

	err = merr.Error(resp)
	if err != nil {
		return err
	}

	return nil
}

func (mgr *channelsMgrImpl) Insert(ctx context.Context, channelMsgs map[string][]*msgpb.InsertRequest) error {
	for channel, msgs := range channelMsgs {
		err := retry.Do(ctx, func() error {
			err := mgr.insert(ctx, channel, msgs)
			if err != nil {
				mgr.UpdateNodeInfo(ctx)
				return err
			}
			return nil
		}, retry.Sleep(200*time.Millisecond))

		if err != nil {
			log.Warn("some channel could not work, some data not insert success", zap.String("channel", channel), zap.Error(err))
			return err
		}
	}

	return nil
}

// implementation assertion
var _ channelsMgr = (*channelsMgrImpl)(nil)

func incPChansMetrics(pchans []pChan) {
	for _, pc := range pchans {
		metrics.ProxyMsgStreamObjectsForPChan.WithLabelValues(strconv.FormatInt(paramtable.GetNodeID(), 10), pc).Inc()
	}
}

func decPChanMetrics(pchans []pChan) {
	for _, pc := range pchans {
		metrics.ProxyMsgStreamObjectsForPChan.WithLabelValues(strconv.FormatInt(paramtable.GetNodeID(), 10), pc).Dec()
	}
}
