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

	"go.uber.org/zap"

	"github.com/milvus-io/milvus-proto/go-api/v2/commonpb"
	"github.com/milvus-io/milvus-proto/go-api/v2/milvuspb"
	"github.com/milvus-io/milvus/internal/distributed/streaming"
	"github.com/milvus-io/milvus/pkg/v2/log"
	"github.com/milvus-io/milvus/pkg/v2/proto/messagespb"
	"github.com/milvus-io/milvus/pkg/v2/streaming/util/message"
	"github.com/milvus-io/milvus/pkg/v2/util/merr"
)

type batchUpdateManifestTask struct {
	baseTask
	Condition
	req          *milvuspb.BatchUpdateManifestRequest
	ctx          context.Context
	node         *Proxy
	msgID        UniqueID
	taskTS       Timestamp
	collectionID UniqueID
	vchannels    []string
}

func (bt *batchUpdateManifestTask) TraceCtx() context.Context {
	return bt.ctx
}

func (bt *batchUpdateManifestTask) ID() UniqueID {
	return bt.msgID
}

func (bt *batchUpdateManifestTask) SetID(uid UniqueID) {
	bt.msgID = uid
}

func (bt *batchUpdateManifestTask) Name() string {
	return "BatchUpdateManifestTask"
}

func (bt *batchUpdateManifestTask) Type() commonpb.MsgType {
	return 0
}

func (bt *batchUpdateManifestTask) BeginTs() Timestamp {
	return bt.taskTS
}

func (bt *batchUpdateManifestTask) EndTs() Timestamp {
	return bt.taskTS
}

func (bt *batchUpdateManifestTask) SetTs(ts Timestamp) {
	bt.taskTS = ts
}

func (bt *batchUpdateManifestTask) OnEnqueue() error {
	return nil
}

func (bt *batchUpdateManifestTask) PreExecute(ctx context.Context) error {
	req := bt.req
	if req.GetCollectionName() == "" {
		return merr.WrapErrParameterInvalidMsg("collection name is empty")
	}
	if len(req.GetItems()) == 0 {
		return merr.WrapErrParameterInvalidMsg("items is empty")
	}

	collectionID, err := globalMetaCache.GetCollectionID(ctx, req.GetDbName(), req.GetCollectionName())
	if err != nil {
		return err
	}
	bt.collectionID = collectionID

	channels, err := bt.node.chMgr.getVChannels(collectionID)
	if err != nil {
		return err
	}
	bt.vchannels = channels
	return nil
}

func (bt *batchUpdateManifestTask) setChannels() error {
	return nil
}

func (bt *batchUpdateManifestTask) getChannels() []pChan {
	return nil
}

func (bt *batchUpdateManifestTask) Execute(ctx context.Context) error {
	items := make([]*messagespb.BatchUpdateManifestItem, 0, len(bt.req.GetItems()))
	for _, item := range bt.req.GetItems() {
		items = append(items, &messagespb.BatchUpdateManifestItem{
			SegmentId:       item.GetSegmentId(),
			ManifestVersion: item.GetManifestVersion(),
		})
	}

	msg, err := message.NewBatchUpdateManifestMessageBuilderV2().
		WithHeader(&message.BatchUpdateManifestMessageHeader{
			CollectionId: bt.collectionID,
		}).
		WithBody(&message.BatchUpdateManifestMessageBody{
			Items: items,
		}).
		WithBroadcast(bt.vchannels).
		BuildBroadcast()
	if err != nil {
		log.Ctx(ctx).Warn("create batch update manifest message failed", zap.Error(err))
		return err
	}

	resp, err := streaming.WAL().Broadcast().Append(ctx, msg)
	if err != nil {
		log.Ctx(ctx).Warn("broadcast batch update manifest msg failed", zap.Error(err))
		return err
	}
	log.Ctx(ctx).Info(
		"broadcast batch update manifest msg success",
		zap.Int64("collectionID", bt.collectionID),
		zap.Int("itemCount", len(items)),
		zap.Uint64("broadcastID", resp.BroadcastID),
	)
	return nil
}

func (bt *batchUpdateManifestTask) PostExecute(ctx context.Context) error {
	return nil
}
