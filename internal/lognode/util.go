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

package lognode

import (
	"github.com/golang/protobuf/proto"
	"github.com/milvus-io/milvus-proto/go-api/v2/commonpb"
	"github.com/milvus-io/milvus/internal/proto/logpb"
	"github.com/milvus-io/milvus/pkg/mq/msgstream"
	"github.com/milvus-io/milvus/pkg/util/merr"
	"github.com/milvus-io/milvus/pkg/util/paramtable"
)

func CheckTargetID[R interface{ GetBase() *commonpb.MsgBase }](req R) bool {
	return req.GetBase().GetTargetID() == paramtable.GetNodeID()
}

func ParseTsMsg(payload []byte, msgType logpb.MessageType) (msgstream.TsMsg, error) {
	switch msgType {
	case logpb.MessageType_INSERT:
		msg := &msgstream.InsertMsg{
			BaseMsg: msgstream.BaseMsg{},
		}
		err := proto.Unmarshal(payload, &msg.InsertRequest)
		if err != nil {
			return nil, err
		}
		msg.BaseMsg.HashValues = make([]uint32, len(msg.InsertRequest.GetRowIDs()))
		return msg, err
	default:
		return nil, merr.WrapErrParameterInvalid("msgtype", "undefineType")
	}
}
