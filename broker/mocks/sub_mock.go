// Code generated by MockGen. DO NOT EDIT.
// Source: my.pkg/race/broker (interfaces: Sub)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSub is a mock of Sub interface
type MockSub struct {
	ctrl     *gomock.Controller
	recorder *MockSubMockRecorder
}

// MockSubMockRecorder is the mock recorder for MockSub
type MockSubMockRecorder struct {
	mock *MockSub
}

// NewMockSub creates a new mock instance
func NewMockSub(ctrl *gomock.Controller) *MockSub {
	mock := &MockSub{ctrl: ctrl}
	mock.recorder = &MockSubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSub) EXPECT() *MockSubMockRecorder {
	return m.recorder
}

// C mocks base method
func (m *MockSub) C() chan<- interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "C")
	ret0, _ := ret[0].(chan<- interface{})
	return ret0
}

// C indicates an expected call of C
func (mr *MockSubMockRecorder) C() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "C", reflect.TypeOf((*MockSub)(nil).C))
}

// Done mocks base method
func (m *MockSub) Done() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Done")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Done indicates an expected call of Done
func (mr *MockSubMockRecorder) Done() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockSub)(nil).Done))
}
