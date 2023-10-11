// Code generated by MockGen. DO NOT EDIT.
// Source: event_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	services "github.com/DIMO-Network/device-data-api/internal/services"
	"go.uber.org/mock/gomock"
)

// MockEventService is a mock of EventService interface.
type MockEventService struct {
	ctrl     *gomock.Controller
	recorder *MockEventServiceMockRecorder
}

// MockEventServiceMockRecorder is the mock recorder for MockEventService.
type MockEventServiceMockRecorder struct {
	mock *MockEventService
}

// NewMockEventService creates a new mock instance.
func NewMockEventService(ctrl *gomock.Controller) *MockEventService {
	mock := &MockEventService{ctrl: ctrl}
	mock.recorder = &MockEventServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventService) EXPECT() *MockEventServiceMockRecorder {
	return m.recorder
}

// Emit mocks base method.
func (m *MockEventService) Emit(event *services.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Emit", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Emit indicates an expected call of Emit.
func (mr *MockEventServiceMockRecorder) Emit(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Emit", reflect.TypeOf((*MockEventService)(nil).Emit), event)
}
