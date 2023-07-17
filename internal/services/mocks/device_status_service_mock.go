// Code generated by MockGen. DO NOT EDIT.
// Source: device_status_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	response "github.com/DIMO-Network/device-data-api/internal/response"
	models "github.com/DIMO-Network/device-data-api/models"
	gomock "github.com/golang/mock/gomock"
)

// MockDeviceStatusService is a mock of DeviceStatusService interface.
type MockDeviceStatusService struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceStatusServiceMockRecorder
}

// MockDeviceStatusServiceMockRecorder is the mock recorder for MockDeviceStatusService.
type MockDeviceStatusServiceMockRecorder struct {
	mock *MockDeviceStatusService
}

// NewMockDeviceStatusService creates a new mock instance.
func NewMockDeviceStatusService(ctrl *gomock.Controller) *MockDeviceStatusService {
	mock := &MockDeviceStatusService{ctrl: ctrl}
	mock.recorder = &MockDeviceStatusServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceStatusService) EXPECT() *MockDeviceStatusServiceMockRecorder {
	return m.recorder
}

// CalculateRange mocks base method.
func (m *MockDeviceStatusService) CalculateRange(ctx context.Context, deviceDefinitionID string, deviceStyleID *string, fuelPercentRemaining float64) (*float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateRange", ctx, deviceDefinitionID, deviceStyleID, fuelPercentRemaining)
	ret0, _ := ret[0].(*float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateRange indicates an expected call of CalculateRange.
func (mr *MockDeviceStatusServiceMockRecorder) CalculateRange(ctx, deviceDefinitionID, deviceStyleID, fuelPercentRemaining interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateRange", reflect.TypeOf((*MockDeviceStatusService)(nil).CalculateRange), ctx, deviceDefinitionID, deviceStyleID, fuelPercentRemaining)
}

// PrepareDeviceStatusInformation mocks base method.
func (m *MockDeviceStatusService) PrepareDeviceStatusInformation(ctx context.Context, deviceData models.UserDeviceDatumSlice, deviceDefinitionID string, deviceStyleID *string, privilegeIDs []int64) response.DeviceSnapshot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareDeviceStatusInformation", ctx, deviceData, deviceDefinitionID, deviceStyleID, privilegeIDs)
	ret0, _ := ret[0].(response.DeviceSnapshot)
	return ret0
}

// PrepareDeviceStatusInformation indicates an expected call of PrepareDeviceStatusInformation.
func (mr *MockDeviceStatusServiceMockRecorder) PrepareDeviceStatusInformation(ctx, deviceData, deviceDefinitionID, deviceStyleID, privilegeIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareDeviceStatusInformation", reflect.TypeOf((*MockDeviceStatusService)(nil).PrepareDeviceStatusInformation), ctx, deviceData, deviceDefinitionID, deviceStyleID, privilegeIDs)
}