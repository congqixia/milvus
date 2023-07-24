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

package querynode

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/milvus-io/milvus/internal/mq/msgstream"
	"github.com/milvus-io/milvus/internal/proto/internalpb"
)

func TestQueryNodeFlowGraph_consumerFlowGraph(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tSafe := newTSafeReplica()

	streamingReplica, err := genSimpleReplica()
	assert.NoError(t, err)

	fac := genFactory()

	fg, err := newQueryNodeFlowGraph(ctx,
		defaultCollectionID,
		streamingReplica,
		tSafe,
		defaultDMLChannel,
		fac)
	assert.NoError(t, err)

	err = fg.consumeFlowGraph(defaultDMLChannel, defaultSubName)
	assert.NoError(t, err)

	fg.close()
}

func TestQueryNodeFlowGraph_seekQueryNodeFlowGraph(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	streamingReplica, err := genSimpleReplica()
	assert.NoError(t, err)

	fac := genFactory()

	tSafe := newTSafeReplica()

	fg, err := newQueryNodeFlowGraph(ctx,
		defaultCollectionID,
		streamingReplica,
		tSafe,
		defaultDMLChannel,
		fac)
	assert.NoError(t, err)

	position := &internalpb.MsgPosition{
		ChannelName: defaultDMLChannel,
		MsgID:       []byte{},
		MsgGroup:    defaultSubName,
		Timestamp:   0,
	}
	err = fg.consumeFlowGraphFromPosition(position)
	assert.Error(t, err)

	fg.close()
}

type DeltaFlowGraphSuite struct {
	suite.Suite

	replica ReplicaInterface
	tsafe   TSafeReplicaInterface
}

func (s *DeltaFlowGraphSuite) SetupTest() {
	var err error
	s.replica, err = genSimpleReplica()
	s.Require().NoError(err)

	s.tsafe = newTSafeReplica()
}

func (s *DeltaFlowGraphSuite) TestNewFgFailed() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockFactory := msgstream.NewMockMqFactory()
	mockFactory.NewMsgStreamFunc = func(ctx context.Context) (msgstream.MsgStream, error) {
		return nil, errors.New("mock error")
	}

	_, err := newQueryNodeDeltaFlowGraph(ctx,
		defaultCollectionID,
		s.replica,
		s.tsafe,
		defaultDeltaChannel,
		mockFactory,
	)

	s.Error(err)
}

func (s *DeltaFlowGraphSuite) TestConsumeFlowGraph() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fac := genFactory()

	fg, err := newQueryNodeDeltaFlowGraph(ctx,
		defaultCollectionID,
		s.replica,
		s.tsafe,
		defaultDeltaChannel,
		fac,
	)
	s.Require().NoError(err)

	defer fg.close()

	err = fg.consumeFlowGraph(defaultDeltaChannel, defaultSubName)
	s.NoError(err)
}

func TestDeltaFlowGraph(t *testing.T) {
	suite.Run(t, new(DeltaFlowGraphSuite))
}
