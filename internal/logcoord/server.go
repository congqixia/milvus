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

package logcoord

import (
	"sync"

	"github.com/milvus-io/milvus/internal/logcoord/meta"
	"github.com/milvus-io/milvus/internal/logcoord/session"
	"github.com/milvus-io/milvus/internal/metastore"
	"github.com/milvus-io/milvus/internal/types"
	"github.com/milvus-io/milvus/internal/util/sessionutil"
	"github.com/milvus-io/milvus/pkg/log"
	"github.com/milvus-io/milvus/pkg/mq/msgstream"
)

type Server struct {
	meta           meta.Meta
	rootCoord      types.RootCoordClient
	streamFactory  msgstream.Factory
	sessionManager *session.SessionManager

	initOnce  sync.Once
	startOnce sync.Once
	stopOnce  sync.Once
}

func NewLogCoord(factory msgstream.Factory) *Server {
	return &Server{streamFactory: factory}
}

func (m *Server) Init(etcdSession *sessionutil.Session, catalog metastore.DataCoordCatalog) error {
	var err error
	m.initOnce.Do(func() {
		m.sessionManager = session.NewSessionManager(session.DefaultLogNodeConnector, etcdSession)

		err = m.meta.Init(m.rootCoord)
		if err != nil {
			return
		}

		err = m.sessionManager.Init(m.meta)
		if err != nil {
			return
		}
		log.Info("log coord init success")
	})
	return err
}

func (m *Server) Start() {
	m.startOnce.Do(func() {
		m.sessionManager.Start()
	})
}

func (m *Server) Stop() {
	m.stopOnce.Do(func() {
		m.sessionManager.Stop()
	})
}

func (m *Server) WatchVChannel(channels ...string) error {
	err := m.meta.AddVChannel(channels...)
	if err != nil {
		return err
	}
	return nil
}

func (m *Server) UnwatchVChannel(channels ...string) error {
	err := m.meta.RemoveVChannel(channels...)
	if err != nil {
		return err
	}
	return nil
}
