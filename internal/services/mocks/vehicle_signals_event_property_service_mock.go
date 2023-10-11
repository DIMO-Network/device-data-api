// Code generated by MockGen. DO NOT EDIT.
// Source: vehicle_signals_event_property_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	models "github.com/DIMO-Network/device-data-api/models"
	grpc "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"go.uber.org/mock/gomock"
)

// MockVehicleSignalsEventPropertyService is a mock of VehicleSignalsEventPropertyService interface.
type MockVehicleSignalsEventPropertyService struct {
	ctrl     *gomock.Controller
	recorder *MockVehicleSignalsEventPropertyServiceMockRecorder
}

// MockVehicleSignalsEventPropertyServiceMockRecorder is the mock recorder for MockVehicleSignalsEventPropertyService.
type MockVehicleSignalsEventPropertyServiceMockRecorder struct {
	mock *MockVehicleSignalsEventPropertyService
}

// NewMockVehicleSignalsEventPropertyService creates a new mock instance.
func NewMockVehicleSignalsEventPropertyService(ctrl *gomock.Controller) *MockVehicleSignalsEventPropertyService {
	mock := &MockVehicleSignalsEventPropertyService{ctrl: ctrl}
	mock.recorder = &MockVehicleSignalsEventPropertyServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVehicleSignalsEventPropertyService) EXPECT() *MockVehicleSignalsEventPropertyServiceMockRecorder {
	return m.recorder
}

// GenerateData mocks base method.
func (m *MockVehicleSignalsEventPropertyService) GenerateData(ctx context.Context, dateKey, integrationID string, ud *models.UserDeviceDatum, deviceDefinition *grpc.GetDeviceDefinitionItemResponse, eventAvailableProperties map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateData", ctx, dateKey, integrationID, ud, deviceDefinition, eventAvailableProperties)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateData indicates an expected call of GenerateData.
func (mr *MockVehicleSignalsEventPropertyServiceMockRecorder) GenerateData(ctx, dateKey, integrationID, ud, deviceDefinition, eventAvailableProperties interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateData", reflect.TypeOf((*MockVehicleSignalsEventPropertyService)(nil).GenerateData), ctx, dateKey, integrationID, ud, deviceDefinition, eventAvailableProperties)
}