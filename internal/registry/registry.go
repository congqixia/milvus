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

package registry

import (
	context "context"

	options "github.com/milvus-io/milvus/internal/registry/options"
)

// Registry is the interface for service registry.
type Registry interface {
	RegisterService(ctx context.Context, service string, addr string, opts ...options.RegisterOption) (ServiceEntry, error)
	DeregisterService(ctx context.Context, entry ServiceEntry) error
	GetService(ctx context.Context, service string) (ServiceEntry, error)
	ListServices(ctx context.Context, service string) ([]ServiceEntry, error)
	WatchService(ctx context.Context, service string, opts ...options.WatchOption) (ServiceWatcher[any], error)
}

// Watcher interface is the abstraction for service watcher.
type Watcher interface {
	Watch() <-chan SessionEvent[any]
	Stop()
}
