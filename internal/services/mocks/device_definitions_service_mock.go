// Code generated by MockGen. DO NOT EDIT.
// Source: device_definitions_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	grpc "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	gomock "github.com/golang/mock/gomock"
)

// MockDeviceDefinitionsAPIService is a mock of DeviceDefinitionsAPIService interface.
type MockDeviceDefinitionsAPIService struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceDefinitionsAPIServiceMockRecorder
}

// MockDeviceDefinitionsAPIServiceMockRecorder is the mock recorder for MockDeviceDefinitionsAPIService.
type MockDeviceDefinitionsAPIServiceMockRecorder struct {
	mock *MockDeviceDefinitionsAPIService
}

// NewMockDeviceDefinitionsAPIService creates a new mock instance.
func NewMockDeviceDefinitionsAPIService(ctrl *gomock.Controller) *MockDeviceDefinitionsAPIService {
	mock := &MockDeviceDefinitionsAPIService{ctrl: ctrl}
	mock.recorder = &MockDeviceDefinitionsAPIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceDefinitionsAPIService) EXPECT() *MockDeviceDefinitionsAPIServiceMockRecorder {
	return m.recorder
}

// GetDeviceDefinition mocks base method.
func (m *MockDeviceDefinitionsAPIService) GetDeviceDefinition(ctx context.Context, id string) (*grpc.GetDeviceDefinitionItemResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceDefinition", ctx, id)
	ret0, _ := ret[0].(*grpc.GetDeviceDefinitionItemResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceDefinition indicates an expected call of GetDeviceDefinition.
func (mr *MockDeviceDefinitionsAPIServiceMockRecorder) GetDeviceDefinition(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceDefinition", reflect.TypeOf((*MockDeviceDefinitionsAPIService)(nil).GetDeviceDefinition), ctx, id)
}