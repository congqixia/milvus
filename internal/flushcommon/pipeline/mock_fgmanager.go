// Code generated by mockery v2.46.0. DO NOT EDIT.

package pipeline

import mock "github.com/stretchr/testify/mock"

// MockFlowgraphManager is an autogenerated mock type for the FlowgraphManager type
type MockFlowgraphManager struct {
	mock.Mock
}

type MockFlowgraphManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFlowgraphManager) EXPECT() *MockFlowgraphManager_Expecter {
	return &MockFlowgraphManager_Expecter{mock: &_m.Mock}
}

// AddFlowgraph provides a mock function with given fields: ds
func (_m *MockFlowgraphManager) AddFlowgraph(ds *DataSyncService) {
	_m.Called(ds)
}

// MockFlowgraphManager_AddFlowgraph_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddFlowgraph'
type MockFlowgraphManager_AddFlowgraph_Call struct {
	*mock.Call
}

// AddFlowgraph is a helper method to define mock.On call
//   - ds *DataSyncService
func (_e *MockFlowgraphManager_Expecter) AddFlowgraph(ds interface{}) *MockFlowgraphManager_AddFlowgraph_Call {
	return &MockFlowgraphManager_AddFlowgraph_Call{Call: _e.mock.On("AddFlowgraph", ds)}
}

func (_c *MockFlowgraphManager_AddFlowgraph_Call) Run(run func(ds *DataSyncService)) *MockFlowgraphManager_AddFlowgraph_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*DataSyncService))
	})
	return _c
}

func (_c *MockFlowgraphManager_AddFlowgraph_Call) Return() *MockFlowgraphManager_AddFlowgraph_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockFlowgraphManager_AddFlowgraph_Call) RunAndReturn(run func(*DataSyncService)) *MockFlowgraphManager_AddFlowgraph_Call {
	_c.Call.Return(run)
	return _c
}

// ClearFlowgraphs provides a mock function with given fields:
func (_m *MockFlowgraphManager) ClearFlowgraphs() {
	_m.Called()
}

// MockFlowgraphManager_ClearFlowgraphs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ClearFlowgraphs'
type MockFlowgraphManager_ClearFlowgraphs_Call struct {
	*mock.Call
}

// ClearFlowgraphs is a helper method to define mock.On call
func (_e *MockFlowgraphManager_Expecter) ClearFlowgraphs() *MockFlowgraphManager_ClearFlowgraphs_Call {
	return &MockFlowgraphManager_ClearFlowgraphs_Call{Call: _e.mock.On("ClearFlowgraphs")}
}

func (_c *MockFlowgraphManager_ClearFlowgraphs_Call) Run(run func()) *MockFlowgraphManager_ClearFlowgraphs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFlowgraphManager_ClearFlowgraphs_Call) Return() *MockFlowgraphManager_ClearFlowgraphs_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockFlowgraphManager_ClearFlowgraphs_Call) RunAndReturn(run func()) *MockFlowgraphManager_ClearFlowgraphs_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields:
func (_m *MockFlowgraphManager) Close() {
	_m.Called()
}

// MockFlowgraphManager_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockFlowgraphManager_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockFlowgraphManager_Expecter) Close() *MockFlowgraphManager_Close_Call {
	return &MockFlowgraphManager_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockFlowgraphManager_Close_Call) Run(run func()) *MockFlowgraphManager_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFlowgraphManager_Close_Call) Return() *MockFlowgraphManager_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockFlowgraphManager_Close_Call) RunAndReturn(run func()) *MockFlowgraphManager_Close_Call {
	_c.Call.Return(run)
	return _c
}

// GetChannelsJSON provides a mock function with given fields: collectionID
func (_m *MockFlowgraphManager) GetChannelsJSON(collectionID int64) string {
	ret := _m.Called(collectionID)

	if len(ret) == 0 {
		panic("no return value specified for GetChannelsJSON")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(int64) string); ok {
		r0 = rf(collectionID)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockFlowgraphManager_GetChannelsJSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChannelsJSON'
type MockFlowgraphManager_GetChannelsJSON_Call struct {
	*mock.Call
}

// GetChannelsJSON is a helper method to define mock.On call
//   - collectionID int64
func (_e *MockFlowgraphManager_Expecter) GetChannelsJSON(collectionID interface{}) *MockFlowgraphManager_GetChannelsJSON_Call {
	return &MockFlowgraphManager_GetChannelsJSON_Call{Call: _e.mock.On("GetChannelsJSON", collectionID)}
}

func (_c *MockFlowgraphManager_GetChannelsJSON_Call) Run(run func(collectionID int64)) *MockFlowgraphManager_GetChannelsJSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockFlowgraphManager_GetChannelsJSON_Call) Return(_a0 string) *MockFlowgraphManager_GetChannelsJSON_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFlowgraphManager_GetChannelsJSON_Call) RunAndReturn(run func(int64) string) *MockFlowgraphManager_GetChannelsJSON_Call {
	_c.Call.Return(run)
	return _c
}

// GetCollectionIDs provides a mock function with given fields:
func (_m *MockFlowgraphManager) GetCollectionIDs() []int64 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetCollectionIDs")
	}

	var r0 []int64
	if rf, ok := ret.Get(0).(func() []int64); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int64)
		}
	}

	return r0
}

// MockFlowgraphManager_GetCollectionIDs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCollectionIDs'
type MockFlowgraphManager_GetCollectionIDs_Call struct {
	*mock.Call
}

// GetCollectionIDs is a helper method to define mock.On call
func (_e *MockFlowgraphManager_Expecter) GetCollectionIDs() *MockFlowgraphManager_GetCollectionIDs_Call {
	return &MockFlowgraphManager_GetCollectionIDs_Call{Call: _e.mock.On("GetCollectionIDs")}
}

func (_c *MockFlowgraphManager_GetCollectionIDs_Call) Run(run func()) *MockFlowgraphManager_GetCollectionIDs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFlowgraphManager_GetCollectionIDs_Call) Return(_a0 []int64) *MockFlowgraphManager_GetCollectionIDs_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFlowgraphManager_GetCollectionIDs_Call) RunAndReturn(run func() []int64) *MockFlowgraphManager_GetCollectionIDs_Call {
	_c.Call.Return(run)
	return _c
}

// GetFlowgraphCount provides a mock function with given fields:
func (_m *MockFlowgraphManager) GetFlowgraphCount() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetFlowgraphCount")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockFlowgraphManager_GetFlowgraphCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFlowgraphCount'
type MockFlowgraphManager_GetFlowgraphCount_Call struct {
	*mock.Call
}

// GetFlowgraphCount is a helper method to define mock.On call
func (_e *MockFlowgraphManager_Expecter) GetFlowgraphCount() *MockFlowgraphManager_GetFlowgraphCount_Call {
	return &MockFlowgraphManager_GetFlowgraphCount_Call{Call: _e.mock.On("GetFlowgraphCount")}
}

func (_c *MockFlowgraphManager_GetFlowgraphCount_Call) Run(run func()) *MockFlowgraphManager_GetFlowgraphCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFlowgraphManager_GetFlowgraphCount_Call) Return(_a0 int) *MockFlowgraphManager_GetFlowgraphCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFlowgraphManager_GetFlowgraphCount_Call) RunAndReturn(run func() int) *MockFlowgraphManager_GetFlowgraphCount_Call {
	_c.Call.Return(run)
	return _c
}

// GetFlowgraphService provides a mock function with given fields: channel
func (_m *MockFlowgraphManager) GetFlowgraphService(channel string) (*DataSyncService, bool) {
	ret := _m.Called(channel)

	if len(ret) == 0 {
		panic("no return value specified for GetFlowgraphService")
	}

	var r0 *DataSyncService
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (*DataSyncService, bool)); ok {
		return rf(channel)
	}
	if rf, ok := ret.Get(0).(func(string) *DataSyncService); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DataSyncService)
		}
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(channel)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockFlowgraphManager_GetFlowgraphService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFlowgraphService'
type MockFlowgraphManager_GetFlowgraphService_Call struct {
	*mock.Call
}

// GetFlowgraphService is a helper method to define mock.On call
//   - channel string
func (_e *MockFlowgraphManager_Expecter) GetFlowgraphService(channel interface{}) *MockFlowgraphManager_GetFlowgraphService_Call {
	return &MockFlowgraphManager_GetFlowgraphService_Call{Call: _e.mock.On("GetFlowgraphService", channel)}
}

func (_c *MockFlowgraphManager_GetFlowgraphService_Call) Run(run func(channel string)) *MockFlowgraphManager_GetFlowgraphService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockFlowgraphManager_GetFlowgraphService_Call) Return(_a0 *DataSyncService, _a1 bool) *MockFlowgraphManager_GetFlowgraphService_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFlowgraphManager_GetFlowgraphService_Call) RunAndReturn(run func(string) (*DataSyncService, bool)) *MockFlowgraphManager_GetFlowgraphService_Call {
	_c.Call.Return(run)
	return _c
}

// GetSegmentsJSON provides a mock function with given fields: collectionID
func (_m *MockFlowgraphManager) GetSegmentsJSON(collectionID int64) string {
	ret := _m.Called(collectionID)

	if len(ret) == 0 {
		panic("no return value specified for GetSegmentsJSON")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(int64) string); ok {
		r0 = rf(collectionID)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockFlowgraphManager_GetSegmentsJSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSegmentsJSON'
type MockFlowgraphManager_GetSegmentsJSON_Call struct {
	*mock.Call
}

// GetSegmentsJSON is a helper method to define mock.On call
//   - collectionID int64
func (_e *MockFlowgraphManager_Expecter) GetSegmentsJSON(collectionID interface{}) *MockFlowgraphManager_GetSegmentsJSON_Call {
	return &MockFlowgraphManager_GetSegmentsJSON_Call{Call: _e.mock.On("GetSegmentsJSON", collectionID)}
}

func (_c *MockFlowgraphManager_GetSegmentsJSON_Call) Run(run func(collectionID int64)) *MockFlowgraphManager_GetSegmentsJSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockFlowgraphManager_GetSegmentsJSON_Call) Return(_a0 string) *MockFlowgraphManager_GetSegmentsJSON_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFlowgraphManager_GetSegmentsJSON_Call) RunAndReturn(run func(int64) string) *MockFlowgraphManager_GetSegmentsJSON_Call {
	_c.Call.Return(run)
	return _c
}

// HasFlowgraph provides a mock function with given fields: channel
func (_m *MockFlowgraphManager) HasFlowgraph(channel string) bool {
	ret := _m.Called(channel)

	if len(ret) == 0 {
		panic("no return value specified for HasFlowgraph")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(channel)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockFlowgraphManager_HasFlowgraph_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasFlowgraph'
type MockFlowgraphManager_HasFlowgraph_Call struct {
	*mock.Call
}

// HasFlowgraph is a helper method to define mock.On call
//   - channel string
func (_e *MockFlowgraphManager_Expecter) HasFlowgraph(channel interface{}) *MockFlowgraphManager_HasFlowgraph_Call {
	return &MockFlowgraphManager_HasFlowgraph_Call{Call: _e.mock.On("HasFlowgraph", channel)}
}

func (_c *MockFlowgraphManager_HasFlowgraph_Call) Run(run func(channel string)) *MockFlowgraphManager_HasFlowgraph_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockFlowgraphManager_HasFlowgraph_Call) Return(_a0 bool) *MockFlowgraphManager_HasFlowgraph_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFlowgraphManager_HasFlowgraph_Call) RunAndReturn(run func(string) bool) *MockFlowgraphManager_HasFlowgraph_Call {
	_c.Call.Return(run)
	return _c
}

// HasFlowgraphWithOpID provides a mock function with given fields: channel, opID
func (_m *MockFlowgraphManager) HasFlowgraphWithOpID(channel string, opID int64) bool {
	ret := _m.Called(channel, opID)

	if len(ret) == 0 {
		panic("no return value specified for HasFlowgraphWithOpID")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, int64) bool); ok {
		r0 = rf(channel, opID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockFlowgraphManager_HasFlowgraphWithOpID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasFlowgraphWithOpID'
type MockFlowgraphManager_HasFlowgraphWithOpID_Call struct {
	*mock.Call
}

// HasFlowgraphWithOpID is a helper method to define mock.On call
//   - channel string
//   - opID int64
func (_e *MockFlowgraphManager_Expecter) HasFlowgraphWithOpID(channel interface{}, opID interface{}) *MockFlowgraphManager_HasFlowgraphWithOpID_Call {
	return &MockFlowgraphManager_HasFlowgraphWithOpID_Call{Call: _e.mock.On("HasFlowgraphWithOpID", channel, opID)}
}

func (_c *MockFlowgraphManager_HasFlowgraphWithOpID_Call) Run(run func(channel string, opID int64)) *MockFlowgraphManager_HasFlowgraphWithOpID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(int64))
	})
	return _c
}

func (_c *MockFlowgraphManager_HasFlowgraphWithOpID_Call) Return(_a0 bool) *MockFlowgraphManager_HasFlowgraphWithOpID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFlowgraphManager_HasFlowgraphWithOpID_Call) RunAndReturn(run func(string, int64) bool) *MockFlowgraphManager_HasFlowgraphWithOpID_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveFlowgraph provides a mock function with given fields: channel
func (_m *MockFlowgraphManager) RemoveFlowgraph(channel string) {
	_m.Called(channel)
}

// MockFlowgraphManager_RemoveFlowgraph_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveFlowgraph'
type MockFlowgraphManager_RemoveFlowgraph_Call struct {
	*mock.Call
}

// RemoveFlowgraph is a helper method to define mock.On call
//   - channel string
func (_e *MockFlowgraphManager_Expecter) RemoveFlowgraph(channel interface{}) *MockFlowgraphManager_RemoveFlowgraph_Call {
	return &MockFlowgraphManager_RemoveFlowgraph_Call{Call: _e.mock.On("RemoveFlowgraph", channel)}
}

func (_c *MockFlowgraphManager_RemoveFlowgraph_Call) Run(run func(channel string)) *MockFlowgraphManager_RemoveFlowgraph_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockFlowgraphManager_RemoveFlowgraph_Call) Return() *MockFlowgraphManager_RemoveFlowgraph_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockFlowgraphManager_RemoveFlowgraph_Call) RunAndReturn(run func(string)) *MockFlowgraphManager_RemoveFlowgraph_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFlowgraphManager creates a new instance of MockFlowgraphManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFlowgraphManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFlowgraphManager {
	mock := &MockFlowgraphManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
