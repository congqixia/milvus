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

package session

import (
	"context"
	"sync"
)

type SessionManager struct {
	sessions struct {
		sync.RWMutex
		data map[int64]*Session
	}
	creator SessionCreator
}

func (s *SessionManager) AddSession(ctx context.Context, nodeID int64, address string) error {
	s.sessions.Lock()
	defer s.sessions.Unlock()

	session := NewSession(nodeID, address, s.creator)
	err := session.Init(ctx)
	if err != nil {
		return err
	}

	s.sessions.data[nodeID] = session
	return nil
}

func (s *SessionManager) GetSessions(nodeID int64) *Session {
	s.sessions.RLock()
	defer s.sessions.RUnlock()

	session, ok := s.sessions.data[nodeID]
	if !ok {
		return nil
	}
	return session
}

func NewSessionManager(creator SessionCreator) *SessionManager {
	return &SessionManager{
		sessions: struct {
			sync.RWMutex
			data map[int64]*Session
		}{data: make(map[int64]*Session)},
		creator: creator,
	}
}
