// Code generated by MockGen. DO NOT EDIT.
// Source: autopi_api_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAutoPiAPIService is a mock of AutoPiAPIService interface.
type MockAutoPiAPIService struct {
	ctrl     *gomock.Controller
	recorder *MockAutoPiAPIServiceMockRecorder
}

// MockAutoPiAPIServiceMockRecorder is the mock recorder for MockAutoPiAPIService.
type MockAutoPiAPIServiceMockRecorder struct {
	mock *MockAutoPiAPIService
}

// NewMockAutoPiAPIService creates a new mock instance.
func NewMockAutoPiAPIService(ctrl *gomock.Controller) *MockAutoPiAPIService {
	mock := &MockAutoPiAPIService{ctrl: ctrl}
	mock.recorder = &MockAutoPiAPIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAutoPiAPIService) EXPECT() *MockAutoPiAPIServiceMockRecorder {
	return m.recorder
}

// UpdateState mocks base method.
func (m *MockAutoPiAPIService) UpdateState(deviceID, state string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateState", deviceID, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateState indicates an expected call of UpdateState.
func (mr *MockAutoPiAPIServiceMockRecorder) UpdateState(deviceID, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateState", reflect.TypeOf((*MockAutoPiAPIService)(nil).UpdateState), deviceID, state)
}
