// Code generated by go generate; DO NOT EDIT
// This file is generated by go generate
package inproc

import (
	"github.com/milvus-io/milvus/internal/types"
	"github.com/milvus-io/milvus/pkg/util/typeutil"
)

type wRootCoord struct {
	types.RootCoord
}

func (w *wRootCoord) Set(v types.RootCoord) {
	w.RootCoord = v
}

func wrapRootCoord(s *inProcServiceDiscovery, v types.RootCoord) types.RootCoord {
	s.serviceMut.Lock()
	defer s.serviceMut.Unlock()
	w := &wRootCoord{RootCoord: v}

	s.wrappers[typeutil.RootCoordRole] = append(s.wrappers[typeutil.RootCoordRole], w)
	return w
}

type wDataCoord struct {
	types.DataCoord
}

func (w *wDataCoord) Set(v types.DataCoord) {
	w.DataCoord = v
}

func wrapDataCoord(s *inProcServiceDiscovery, v types.DataCoord) types.DataCoord {
	s.serviceMut.Lock()
	defer s.serviceMut.Unlock()
	w := &wDataCoord{DataCoord: v}

	s.wrappers[typeutil.DataCoordRole] = append(s.wrappers[typeutil.DataCoordRole], w)
	return w
}

type wQueryCoord struct {
	types.QueryCoord
}

func (w *wQueryCoord) Set(v types.QueryCoord) {
	w.QueryCoord = v
}

func wrapQueryCoord(s *inProcServiceDiscovery, v types.QueryCoord) types.QueryCoord {
	s.serviceMut.Lock()
	defer s.serviceMut.Unlock()
	w := &wQueryCoord{QueryCoord: v}

	s.wrappers[typeutil.QueryCoordRole] = append(s.wrappers[typeutil.QueryCoordRole], w)
	return w
}
