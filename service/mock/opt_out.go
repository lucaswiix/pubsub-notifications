// Code generated by MockGen. DO NOT EDIT.
// Source: meli/notifications/service (interfaces: OptOutService)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOptOutService is a mock of OptOutService interface.
type MockOptOutService struct {
	ctrl     *gomock.Controller
	recorder *MockOptOutServiceMockRecorder
}

// MockOptOutServiceMockRecorder is the mock recorder for MockOptOutService.
type MockOptOutServiceMockRecorder struct {
	mock *MockOptOutService
}

// NewMockOptOutService creates a new mock instance.
func NewMockOptOutService(ctrl *gomock.Controller) *MockOptOutService {
	mock := &MockOptOutService{ctrl: ctrl}
	mock.recorder = &MockOptOutServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOptOutService) EXPECT() *MockOptOutServiceMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockOptOutService) Del(arg0 string, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockOptOutServiceMockRecorder) Del(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockOptOutService)(nil).Del), arg0, arg1)
}

// Is mocks base method.
func (m *MockOptOutService) Is(arg0 string, arg1 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Is", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Is indicates an expected call of Is.
func (mr *MockOptOutServiceMockRecorder) Is(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Is", reflect.TypeOf((*MockOptOutService)(nil).Is), arg0, arg1)
}

// Set mocks base method.
func (m *MockOptOutService) Set(arg0 string, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockOptOutServiceMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockOptOutService)(nil).Set), arg0, arg1)
}
