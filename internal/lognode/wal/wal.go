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

package wal

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/milvus-io/milvus/pkg/log"
	"github.com/milvus-io/milvus/pkg/mq/msgstream"
	"github.com/milvus-io/milvus/pkg/util/timerecord"
	"go.uber.org/zap"
)

type WriteAheadLogger struct {
	channel string
	lastTs  atomic.Uint64

	stream msgstream.MsgStream
	mu     sync.Mutex
}

func NewWriteAheadLogger(channel string) *WriteAheadLogger {
	return &WriteAheadLogger{
		channel: channel,
	}
}

func (logger *WriteAheadLogger) GetChannelName() string {
	return logger.channel
}

func (logger *WriteAheadLogger) GetLastTimestamp() uint64 {
	return logger.lastTs.Load()
}

func (logger *WriteAheadLogger) Init(ctx context.Context, factory msgstream.Factory) error {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	stream, err := factory.NewMsgStream(ctx)
	if err != nil {
		return err
	}

	logger.stream = stream
	logger.stream.AsProducer([]string{logger.channel})
	return nil
}

func (logger *WriteAheadLogger) Produce(ctx context.Context, pack *msgstream.MsgPack) error {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	record := timerecord.NewTimeRecorder("produce")
	err := logger.stream.Produce(pack)
	if err != nil {
		return err
	}
	log.Info("test produce", zap.Duration("interval", record.ElapseSpan()))

	logger.lastTs.Store(pack.EndTs)
	return nil
}
