// Code generated by MockGen. DO NOT EDIT.
// Source: meli/notifications/repository (interfaces: OptOutRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOptOutRepository is a mock of OptOutRepository interface.
type MockOptOutRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOptOutRepositoryMockRecorder
}

// MockOptOutRepositoryMockRecorder is the mock recorder for MockOptOutRepository.
type MockOptOutRepositoryMockRecorder struct {
	mock *MockOptOutRepository
}

// NewMockOptOutRepository creates a new mock instance.
func NewMockOptOutRepository(ctrl *gomock.Controller) *MockOptOutRepository {
	mock := &MockOptOutRepository{ctrl: ctrl}
	mock.recorder = &MockOptOutRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOptOutRepository) EXPECT() *MockOptOutRepositoryMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockOptOutRepository) Del(arg0 string, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockOptOutRepositoryMockRecorder) Del(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockOptOutRepository)(nil).Del), arg0, arg1)
}

// Is mocks base method.
func (m *MockOptOutRepository) Is(arg0 string, arg1 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Is", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Is indicates an expected call of Is.
func (mr *MockOptOutRepositoryMockRecorder) Is(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Is", reflect.TypeOf((*MockOptOutRepository)(nil).Is), arg0, arg1)
}

// Set mocks base method.
func (m *MockOptOutRepository) Set(arg0 string, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockOptOutRepositoryMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockOptOutRepository)(nil).Set), arg0, arg1)
}
