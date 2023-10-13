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
	"fmt"
	"strings"

	"github.com/milvus-io/milvus/pkg/util/paramtable"
)

func genPChannelNames(prefix string, num int) []string {
	var results []string
	for idx := 0; idx < num; idx++ {
		result := fmt.Sprintf("%s_%d", prefix, idx)
		results = append(results, result)
	}
	return results
}

func getPChannelList() []string {
	params := paramtable.Get()

	if params.CommonCfg.PreCreatedTopicEnabled.GetAsBool() {
		return params.CommonCfg.TopicNames.GetAsStrings()
	} else {
		prefix := params.CommonCfg.RootCoordDml.GetValue()
		num := params.RootCoordCfg.DmlChannelNum.GetAsInt()
		return genPChannelNames(prefix, num)
	}
}

func getPChannelName(vchannel string) string {
	index := strings.LastIndex(vchannel, "_")
	if index < 0 {
		return vchannel
	}
	return vchannel[:index]
}
