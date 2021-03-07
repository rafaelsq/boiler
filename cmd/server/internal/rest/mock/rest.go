// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock is a generated GoMock package.
package mock

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockResp is a mock of Resp interface.
type MockResp struct {
	ctrl     *gomock.Controller
	recorder *MockRespMockRecorder
}

// MockRespMockRecorder is the mock recorder for MockResp.
type MockRespMockRecorder struct {
	mock *MockResp
}

// NewMockResp creates a new mock instance.
func NewMockResp(ctrl *gomock.Controller) *MockResp {
	mock := &MockResp{ctrl: ctrl}
	mock.recorder = &MockRespMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResp) EXPECT() *MockRespMockRecorder {
	return m.recorder
}

// Fail mocks base method.
func (m *MockResp) Fail(w http.ResponseWriter, r *http.Request, err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Fail", w, r, err)
}

// Fail indicates an expected call of Fail.
func (mr *MockRespMockRecorder) Fail(w, r, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fail", reflect.TypeOf((*MockResp)(nil).Fail), w, r, err)
}

// Failf mocks base method.
func (m *MockResp) Failf(w http.ResponseWriter, r *http.Request, format string, a ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{w, r, format}
	for _, a_2 := range a {
		varargs = append(varargs, a_2)
	}
	m.ctrl.Call(m, "Failf", varargs...)
}

// Failf indicates an expected call of Failf.
func (mr *MockRespMockRecorder) Failf(w, r, format interface{}, a ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{w, r, format}, a...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Failf", reflect.TypeOf((*MockResp)(nil).Failf), varargs...)
}

// JSON mocks base method.
func (m *MockResp) JSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "JSON", w, r, data)
}

// JSON indicates an expected call of JSON.
func (mr *MockRespMockRecorder) JSON(w, r, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JSON", reflect.TypeOf((*MockResp)(nil).JSON), w, r, data)
}
