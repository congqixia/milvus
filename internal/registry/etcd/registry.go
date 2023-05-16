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

package etcd

import (
	"context"

	"github.com/milvus-io/milvus/internal/registry"
	"github.com/milvus-io/milvus/internal/registry/common"
	"github.com/milvus-io/milvus/internal/registry/options"
	"github.com/milvus-io/milvus/pkg/util/merr"
)

type etcdRegistry struct {
}

func (e *etcdRegistry) RegisterService(ctx context.Context, service string, entry registry.ServiceEntry, opts ...options.RegisterOption) error {
	var exclusive bool
	switch {
	case common.IsCoordinator(service):
		exclusive = true
	case common.IsNode(service):
		exclusive = false
	default:
		return merr.WrapErrParameterInvalid("legal component type", service)
	}

	opt := options.DefaultSessionOpt()
	for _, o := range opts {
		o(&opt)
	}
	opt.Exclusive = exclusive

	return nil
}

func (r *etcdRegistry) DeregisterService(ctx context.Context, service string, entry registry.ServiceEntry) error {
	panic("not implemented") // TODO: Implement
}

func (r *etcdRegistry) GetService(ctx context.Context, service string) (registry.ServiceEntry, error) {
	panic("not implemented") // TODO: Implement
}

func (r *etcdRegistry) ListServices(ctx context.Context, service string) ([]registry.ServiceEntry, error) {
	panic("not implemented") // TODO: Implement
}

func (r *etcdRegistry) WatchService(ctx context.Context, service string, opts ...options.WatchOption) (registry.ServiceWatcher[any], error) {
	panic("not implemented") // TODO: Implement
}
