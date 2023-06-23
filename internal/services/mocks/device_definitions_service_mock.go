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

// GetDeviceDefinitionByID mocks base method.
func (m *MockDeviceDefinitionsAPIService) GetDeviceDefinitionByID(ctx context.Context, id string) (*grpc.GetDeviceDefinitionItemResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceDefinitionByID", ctx, id)
	ret0, _ := ret[0].(*grpc.GetDeviceDefinitionItemResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceDefinitionByID indicates an expected call of GetDeviceDefinitionByID.
func (mr *MockDeviceDefinitionsAPIServiceMockRecorder) GetDeviceDefinitionByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceDefinitionByID", reflect.TypeOf((*MockDeviceDefinitionsAPIService)(nil).GetDeviceDefinitionByID), ctx, id)
}

// GetDeviceDefinitionsByIDs mocks base method.
func (m *MockDeviceDefinitionsAPIService) GetDeviceDefinitionsByIDs(ctx context.Context, ids []string) ([]*grpc.GetDeviceDefinitionItemResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceDefinitionsByIDs", ctx, ids)
	ret0, _ := ret[0].([]*grpc.GetDeviceDefinitionItemResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceDefinitionsByIDs indicates an expected call of GetDeviceDefinitionsByIDs.
func (mr *MockDeviceDefinitionsAPIServiceMockRecorder) GetDeviceDefinitionsByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceDefinitionsByIDs", reflect.TypeOf((*MockDeviceDefinitionsAPIService)(nil).GetDeviceDefinitionsByIDs), ctx, ids)
}

// GetIntegrations mocks base method.
func (m *MockDeviceDefinitionsAPIService) GetIntegrations(ctx context.Context) ([]*grpc.Integration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIntegrations", ctx)
	ret0, _ := ret[0].([]*grpc.Integration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIntegrations indicates an expected call of GetIntegrations.
func (mr *MockDeviceDefinitionsAPIServiceMockRecorder) GetIntegrations(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIntegrations", reflect.TypeOf((*MockDeviceDefinitionsAPIService)(nil).GetIntegrations), ctx)
}
