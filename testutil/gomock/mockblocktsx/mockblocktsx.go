// Code generated by MockGen. DO NOT EDIT.
// Source: ethereum-proxy/internal/blocktsx (interfaces: BlockProvider)

// Package mockblocktsx is a generated GoMock package.
package mockblocktsx

import (
	application "github.com/miscoler/ethereum-proxy/internal/application"
	blocktsx "github.com/miscoler/ethereum-proxy/internal/blocktsx"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBlockProvider is a mock of BlockProvider interface
type MockBlockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockBlockProviderMockRecorder
}

// MockBlockProviderMockRecorder is the mock recorder for MockBlockProvider
type MockBlockProviderMockRecorder struct {
	mock *MockBlockProvider
}

// NewMockBlockProvider creates a new mock instance
func NewMockBlockProvider(ctrl *gomock.Controller) *MockBlockProvider {
	mock := &MockBlockProvider{ctrl: ctrl}
	mock.recorder = &MockBlockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockProvider) EXPECT() *MockBlockProviderMockRecorder {
	return m.recorder
}

// GetBlock mocks base method
func (m *MockBlockProvider) GetBlock(arg0 *application.EContext, arg1 int64, arg2 bool) (*blocktsx.BlockStored, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", arg0, arg1, arg2)
	ret0, _ := ret[0].(*blocktsx.BlockStored)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock
func (mr *MockBlockProviderMockRecorder) GetBlock(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockBlockProvider)(nil).GetBlock), arg0, arg1, arg2)
}

// GetTSXbyHash mocks base method
func (m *MockBlockProvider) GetTSXbyHash(arg0 *blocktsx.BlockStored, arg1 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTSXbyHash", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTSXbyHash indicates an expected call of GetTSXbyHash
func (mr *MockBlockProviderMockRecorder) GetTSXbyHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTSXbyHash", reflect.TypeOf((*MockBlockProvider)(nil).GetTSXbyHash), arg0, arg1)
}

// GetTSXbyIndex mocks base method
func (m *MockBlockProvider) GetTSXbyIndex(arg0 *blocktsx.BlockStored, arg1 int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTSXbyIndex", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTSXbyIndex indicates an expected call of GetTSXbyIndex
func (mr *MockBlockProviderMockRecorder) GetTSXbyIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTSXbyIndex", reflect.TypeOf((*MockBlockProvider)(nil).GetTSXbyIndex), arg0, arg1)
}