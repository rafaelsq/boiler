// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock is a generated GoMock package.
package mock

import (
	entity "boiler/pkg/entity"
	iface "boiler/pkg/iface"
	context "context"
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddUser mocks base method
func (m *MockService) AddUser(arg0 context.Context, arg1, arg2 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser
func (mr *MockServiceMockRecorder) AddUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockService)(nil).AddUser), arg0, arg1, arg2)
}

// DeleteUser mocks base method
func (m *MockService) DeleteUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockServiceMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockService)(nil).DeleteUser), arg0, arg1)
}

// FilterUsers mocks base method
func (m *MockService) FilterUsers(arg0 context.Context, arg1 iface.FilterUsers) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterUsers", arg0, arg1)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterUsers indicates an expected call of FilterUsers
func (mr *MockServiceMockRecorder) FilterUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterUsers", reflect.TypeOf((*MockService)(nil).FilterUsers), arg0, arg1)
}

// GetUserByID mocks base method
func (m *MockService) GetUserByID(arg0 context.Context, arg1 int64) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID
func (mr *MockServiceMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockService)(nil).GetUserByID), arg0, arg1)
}

// GetUserByEmail mocks base method
func (m *MockService) GetUserByEmail(arg0 context.Context, arg1 string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail
func (mr *MockServiceMockRecorder) GetUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockService)(nil).GetUserByEmail), arg0, arg1)
}

// AuthUser mocks base method
func (m *MockService) AuthUser(arg0 context.Context, arg1, arg2 string) (*entity.User, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AuthUser indicates an expected call of AuthUser
func (mr *MockServiceMockRecorder) AuthUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthUser", reflect.TypeOf((*MockService)(nil).AuthUser), arg0, arg1, arg2)
}

// FilterEmails mocks base method
func (m *MockService) FilterEmails(arg0 context.Context, arg1 iface.FilterEmails) ([]*entity.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterEmails", arg0, arg1)
	ret0, _ := ret[0].([]*entity.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterEmails indicates an expected call of FilterEmails
func (mr *MockServiceMockRecorder) FilterEmails(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterEmails", reflect.TypeOf((*MockService)(nil).FilterEmails), arg0, arg1)
}

// AddEmail mocks base method
func (m *MockService) AddEmail(arg0 context.Context, arg1 int64, arg2 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEmail", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEmail indicates an expected call of AddEmail
func (mr *MockServiceMockRecorder) AddEmail(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEmail", reflect.TypeOf((*MockService)(nil).AddEmail), arg0, arg1, arg2)
}

// DeleteEmail mocks base method
func (m *MockService) DeleteEmail(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEmail", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEmail indicates an expected call of DeleteEmail
func (mr *MockServiceMockRecorder) DeleteEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEmail", reflect.TypeOf((*MockService)(nil).DeleteEmail), arg0, arg1)
}

// EnqueueDeleteEmail mocks base method
func (m *MockService) EnqueueDeleteEmail(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnqueueDeleteEmail", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnqueueDeleteEmail indicates an expected call of EnqueueDeleteEmail
func (mr *MockServiceMockRecorder) EnqueueDeleteEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnqueueDeleteEmail", reflect.TypeOf((*MockService)(nil).EnqueueDeleteEmail), arg0, arg1)
}

// AuthUserMiddleware mocks base method
func (m *MockService) AuthUserMiddleware(arg0 http.Handler) http.Handler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthUserMiddleware", arg0)
	ret0, _ := ret[0].(http.Handler)
	return ret0
}

// AuthUserMiddleware indicates an expected call of AuthUserMiddleware
func (mr *MockServiceMockRecorder) AuthUserMiddleware(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthUserMiddleware", reflect.TypeOf((*MockService)(nil).AuthUserMiddleware), arg0)
}
