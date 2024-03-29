// Code generated by MockGen. DO NOT EDIT.
// Source: device_api_service.go
//
// Generated by this command:
//
//	mockgen -source device_api_service.go -destination mocks/device_api_service_mock.go
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	grpc "github.com/DIMO-Network/devices-api/pkg/grpc"
	gomock "go.uber.org/mock/gomock"
)

// MockDeviceAPIService is a mock of DeviceAPIService interface.
type MockDeviceAPIService struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceAPIServiceMockRecorder
}

// MockDeviceAPIServiceMockRecorder is the mock recorder for MockDeviceAPIService.
type MockDeviceAPIServiceMockRecorder struct {
	mock *MockDeviceAPIService
}

// NewMockDeviceAPIService creates a new mock instance.
func NewMockDeviceAPIService(ctrl *gomock.Controller) *MockDeviceAPIService {
	mock := &MockDeviceAPIService{ctrl: ctrl}
	mock.recorder = &MockDeviceAPIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceAPIService) EXPECT() *MockDeviceAPIServiceMockRecorder {
	return m.recorder
}

// GetUserDevice mocks base method.
func (m *MockDeviceAPIService) GetUserDevice(ctx context.Context, userDeviceID string) (*grpc.UserDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDevice", ctx, userDeviceID)
	ret0, _ := ret[0].(*grpc.UserDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDevice indicates an expected call of GetUserDevice.
func (mr *MockDeviceAPIServiceMockRecorder) GetUserDevice(ctx, userDeviceID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDevice", reflect.TypeOf((*MockDeviceAPIService)(nil).GetUserDevice), ctx, userDeviceID)
}

// GetUserDeviceByEthAddr mocks base method.
func (m *MockDeviceAPIService) GetUserDeviceByEthAddr(ctx context.Context, ethAddr []byte) (*grpc.UserDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDeviceByEthAddr", ctx, ethAddr)
	ret0, _ := ret[0].(*grpc.UserDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDeviceByEthAddr indicates an expected call of GetUserDeviceByEthAddr.
func (mr *MockDeviceAPIServiceMockRecorder) GetUserDeviceByEthAddr(ctx, ethAddr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDeviceByEthAddr", reflect.TypeOf((*MockDeviceAPIService)(nil).GetUserDeviceByEthAddr), ctx, ethAddr)
}

// GetUserDeviceByTokenID mocks base method.
func (m *MockDeviceAPIService) GetUserDeviceByTokenID(ctx context.Context, tokenID int64) (*grpc.UserDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDeviceByTokenID", ctx, tokenID)
	ret0, _ := ret[0].(*grpc.UserDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDeviceByTokenID indicates an expected call of GetUserDeviceByTokenID.
func (mr *MockDeviceAPIServiceMockRecorder) GetUserDeviceByTokenID(ctx, tokenID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDeviceByTokenID", reflect.TypeOf((*MockDeviceAPIService)(nil).GetUserDeviceByTokenID), ctx, tokenID)
}

// ListUserDevicesForUser mocks base method.
func (m *MockDeviceAPIService) ListUserDevicesForUser(ctx context.Context, userID string) (*grpc.ListUserDevicesForUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserDevicesForUser", ctx, userID)
	ret0, _ := ret[0].(*grpc.ListUserDevicesForUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserDevicesForUser indicates an expected call of ListUserDevicesForUser.
func (mr *MockDeviceAPIServiceMockRecorder) ListUserDevicesForUser(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserDevicesForUser", reflect.TypeOf((*MockDeviceAPIService)(nil).ListUserDevicesForUser), ctx, userID)
}

// UpdateStatus mocks base method.
func (m *MockDeviceAPIService) UpdateStatus(ctx context.Context, userDeviceID, integrationID, status string) (*grpc.UserDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, userDeviceID, integrationID, status)
	ret0, _ := ret[0].(*grpc.UserDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockDeviceAPIServiceMockRecorder) UpdateStatus(ctx, userDeviceID, integrationID, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockDeviceAPIService)(nil).UpdateStatus), ctx, userDeviceID, integrationID, status)
}

// UserDeviceBelongsToUserID mocks base method.
func (m *MockDeviceAPIService) UserDeviceBelongsToUserID(ctx context.Context, userID, userDeviceID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserDeviceBelongsToUserID", ctx, userID, userDeviceID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserDeviceBelongsToUserID indicates an expected call of UserDeviceBelongsToUserID.
func (mr *MockDeviceAPIServiceMockRecorder) UserDeviceBelongsToUserID(ctx, userID, userDeviceID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserDeviceBelongsToUserID", reflect.TypeOf((*MockDeviceAPIService)(nil).UserDeviceBelongsToUserID), ctx, userID, userDeviceID)
}
