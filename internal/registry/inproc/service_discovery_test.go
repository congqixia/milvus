package inproc

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/milvus-io/milvus/internal/mocks"
	"github.com/milvus-io/milvus/internal/registry"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type InProcSDSuite struct {
	inner *registry.MockServiceDiscovery
	sd    registry.ServiceDiscovery
	suite.Suite
}

func (s *InProcSDSuite) SetupSuite() {
}

func (s *InProcSDSuite) SetupTest() {
	s.inner = registry.NewMockServiceDiscovery(s.T())
	s.sd = NewInProcServiceDiscovery(s.inner)
}

func (s *InProcSDSuite) resetMock() {
	if s.inner != nil {
		s.inner.ExpectedCalls = nil
	}
}

func (s *InProcSDSuite) TestGetServices() {

	result := []registry.ServiceEntry{}
	s.inner.EXPECT().GetServices(mock.Anything, mock.AnythingOfType("string")).Return(result, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sr, err := s.sd.GetServices(ctx, "component")
	s.NoError(err)
	s.Equal(result, sr)
}

func (s *InProcSDSuite) TestGetRootCoord() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run("inner_returns_error", func() {
		defer s.resetMock()

		s.inner.EXPECT().GetRootCoord(mock.Anything).Return(nil, errors.New("mocks"))

		_, err := s.sd.GetRootCoord(ctx)
		s.Error(err)
	})

	s.Run("fallthrough_to_inner", func() {
		defer s.resetMock()
		mrc := mocks.NewRootCoord(s.T())
		s.inner.EXPECT().GetRootCoord(mock.Anything).Return(mrc, nil)

		rc, err := s.sd.GetRootCoord(ctx)
		s.NoError(err)

		wrc, ok := rc.(*wRootCoord)
		s.True(ok)
		s.Equal(mrc, wrc.RootCoord)
	})

	s.Run("from_inproc", func() {
		mrc := mocks.NewRootCoord(s.T())
		s.inner.EXPECT().RegisterRootCoord(mock.Anything, mrc, "mock_rootcoord").Return(nil, nil)
		_, err := s.sd.RegisterRootCoord(ctx, mrc, "mock_rootcoord")
		s.NoError(err)

		rc, err := s.sd.GetRootCoord(ctx)
		s.NoError(err)

		s.Equal(mrc, rc)
	})
}

func (s *InProcSDSuite) TestGetQueryCoord() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run("inner_returns_error", func() {
		defer s.resetMock()

		s.inner.EXPECT().GetQueryCoord(mock.Anything).Return(nil, errors.New("mocks"))

		_, err := s.sd.GetQueryCoord(ctx)
		s.Error(err)
	})

	s.Run("fallthrough_to_inner", func() {
		defer s.resetMock()
		mqc := mocks.NewQueryCoord(s.T())
		s.inner.EXPECT().GetQueryCoord(mock.Anything).Return(mqc, nil)

		rc, err := s.sd.GetQueryCoord(ctx)
		s.NoError(err)

		wrc, ok := rc.(*wQueryCoord)
		s.True(ok)
		s.Equal(mqc, wrc.QueryCoord)
	})

	s.Run("from_inproc", func() {
		mrc := mocks.NewQueryCoord(s.T())
		s.inner.EXPECT().RegisterQueryCoord(mock.Anything, mrc, "mock_querycoord").Return(nil, nil)
		_, err := s.sd.RegisterQueryCoord(ctx, mrc, "mock_querycoord")
		s.NoError(err)

		rc, err := s.sd.GetQueryCoord(ctx)
		s.NoError(err)

		s.Equal(mrc, rc)
	})
}

func (s *InProcSDSuite) TestGetDataCoord() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Run("inner_returns_error", func() {
		defer s.resetMock()

		s.inner.EXPECT().GetDataCoord(mock.Anything).Return(nil, errors.New("mocks"))

		_, err := s.sd.GetDataCoord(ctx)
		s.Error(err)
	})

	s.Run("fallthrough_to_inner", func() {
		defer s.resetMock()
		mqc := mocks.NewDataCoord(s.T())
		s.inner.EXPECT().GetDataCoord(mock.Anything).Return(mqc, nil)

		rc, err := s.sd.GetDataCoord(ctx)
		s.NoError(err)

		wrc, ok := rc.(*wDataCoord)
		s.True(ok)
		s.Equal(mqc, wrc.DataCoord)
	})

	s.Run("from_inproc", func() {
		mrc := mocks.NewDataCoord(s.T())
		s.inner.EXPECT().RegisterDataCoord(mock.Anything, mrc, "mock_datacoord").Return(nil, nil)
		_, err := s.sd.RegisterDataCoord(ctx, mrc, "mock_datacoord")
		s.NoError(err)

		rc, err := s.sd.GetDataCoord(ctx)
		s.NoError(err)

		s.Equal(mrc, rc)
	})
}

func (s *InProcSDSuite) TestWatchProxy() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/*proxies, watcher,*/
	_, _, err := s.sd.WatchProxy(ctx)
	s.NoError(err)
}

func TestInProcServiceDiscovery(t *testing.T) {
	suite.Run(t, new(InProcSDSuite))
}
