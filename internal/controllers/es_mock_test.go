// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DIMO-Network/device-data-api/internal/controllers (interfaces: EsInterface)
//
// Generated by this command:
//
//	mockgen -package controllers_test -destination es_mock_test.go github.com/DIMO-Network/device-data-api/internal/controllers EsInterface
//

// Package controllers_test is a generated GoMock package.
package controllers_test

import (
	context "context"
	json "encoding/json"
	reflect "reflect"

	elastic "github.com/DIMO-Network/device-data-api/internal/services/elastic"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	gomock "go.uber.org/mock/gomock"
)

// MockEsInterface is a mock of EsInterface interface.
type MockEsInterface struct {
	ctrl     *gomock.Controller
	recorder *MockEsInterfaceMockRecorder
}

// MockEsInterfaceMockRecorder is the mock recorder for MockEsInterface.
type MockEsInterfaceMockRecorder struct {
	mock *MockEsInterface
}

// NewMockEsInterface creates a new mock instance.
func NewMockEsInterface(ctrl *gomock.Controller) *MockEsInterface {
	mock := &MockEsInterface{ctrl: ctrl}
	mock.recorder = &MockEsInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEsInterface) EXPECT() *MockEsInterfaceMockRecorder {
	return m.recorder
}

// ESClient mocks base method.
func (m *MockEsInterface) ESClient() *elasticsearch.TypedClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ESClient")
	ret0, _ := ret[0].(*elasticsearch.TypedClient)
	return ret0
}

// ESClient indicates an expected call of ESClient.
func (mr *MockEsInterfaceMockRecorder) ESClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ESClient", reflect.TypeOf((*MockEsInterface)(nil).ESClient))
}

// GetHistory mocks base method.
func (m *MockEsInterface) GetHistory(arg0 context.Context, arg1 elastic.GetHistoryParams) ([]json.RawMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistory", arg0, arg1)
	ret0, _ := ret[0].([]json.RawMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHistory indicates an expected call of GetHistory.
func (mr *MockEsInterfaceMockRecorder) GetHistory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistory", reflect.TypeOf((*MockEsInterface)(nil).GetHistory), arg0, arg1)
}

// GetTotalDailyDistanceDriven mocks base method.
func (m *MockEsInterface) GetTotalDailyDistanceDriven(arg0 context.Context, arg1, arg2 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalDailyDistanceDriven", arg0, arg1, arg2)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalDailyDistanceDriven indicates an expected call of GetTotalDailyDistanceDriven.
func (mr *MockEsInterfaceMockRecorder) GetTotalDailyDistanceDriven(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalDailyDistanceDriven", reflect.TypeOf((*MockEsInterface)(nil).GetTotalDailyDistanceDriven), arg0, arg1, arg2)
}

// GetTotalDistanceDriven mocks base method.
func (m *MockEsInterface) GetTotalDistanceDriven(arg0 context.Context, arg1 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalDistanceDriven", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalDistanceDriven indicates an expected call of GetTotalDistanceDriven.
func (mr *MockEsInterfaceMockRecorder) GetTotalDistanceDriven(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalDistanceDriven", reflect.TypeOf((*MockEsInterface)(nil).GetTotalDistanceDriven), arg0, arg1)
}
