// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock is a generated GoMock package.
package mock

import (
	entity "boiler/pkg/entity"
	store "boiler/pkg/store"
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// AddEmail mocks base method.
func (m *MockInterface) AddEmail(ctx context.Context, tx *sql.Tx, email *entity.Email) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEmail", ctx, tx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEmail indicates an expected call of AddEmail.
func (mr *MockInterfaceMockRecorder) AddEmail(ctx, tx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEmail", reflect.TypeOf((*MockInterface)(nil).AddEmail), ctx, tx, email)
}

// AddUser mocks base method.
func (m *MockInterface) AddUser(ctx context.Context, tx *sql.Tx, user *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, tx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockInterfaceMockRecorder) AddUser(ctx, tx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockInterface)(nil).AddUser), ctx, tx, user)
}

// DeleteEmail mocks base method.
func (m *MockInterface) DeleteEmail(ctx context.Context, tx *sql.Tx, email int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEmail", ctx, tx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEmail indicates an expected call of DeleteEmail.
func (mr *MockInterfaceMockRecorder) DeleteEmail(ctx, tx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEmail", reflect.TypeOf((*MockInterface)(nil).DeleteEmail), ctx, tx, email)
}

// DeleteEmailsByUserID mocks base method.
func (m *MockInterface) DeleteEmailsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEmailsByUserID", ctx, tx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEmailsByUserID indicates an expected call of DeleteEmailsByUserID.
func (mr *MockInterfaceMockRecorder) DeleteEmailsByUserID(ctx, tx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEmailsByUserID", reflect.TypeOf((*MockInterface)(nil).DeleteEmailsByUserID), ctx, tx, userID)
}

// DeleteUser mocks base method.
func (m *MockInterface) DeleteUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, tx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockInterfaceMockRecorder) DeleteUser(ctx, tx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockInterface)(nil).DeleteUser), ctx, tx, userID)
}

// FetchUsers mocks base method.
func (m *MockInterface) FetchUsers(ctx context.Context, ID []int64, users *[]entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchUsers", ctx, ID, users)
	ret0, _ := ret[0].(error)
	return ret0
}

// FetchUsers indicates an expected call of FetchUsers.
func (mr *MockInterfaceMockRecorder) FetchUsers(ctx, ID, users interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUsers", reflect.TypeOf((*MockInterface)(nil).FetchUsers), ctx, ID, users)
}

// FilterEmails mocks base method.
func (m *MockInterface) FilterEmails(ctx context.Context, filter store.FilterEmails, emails *[]entity.Email) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterEmails", ctx, filter, emails)
	ret0, _ := ret[0].(error)
	return ret0
}

// FilterEmails indicates an expected call of FilterEmails.
func (mr *MockInterfaceMockRecorder) FilterEmails(ctx, filter, emails interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterEmails", reflect.TypeOf((*MockInterface)(nil).FilterEmails), ctx, filter, emails)
}

// FilterUsersID mocks base method.
func (m *MockInterface) FilterUsersID(ctx context.Context, filter store.FilterUsers, IDs *[]int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterUsersID", ctx, filter, IDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// FilterUsersID indicates an expected call of FilterUsersID.
func (mr *MockInterfaceMockRecorder) FilterUsersID(ctx, filter, IDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterUsersID", reflect.TypeOf((*MockInterface)(nil).FilterUsersID), ctx, filter, IDs)
}

// Tx mocks base method.
func (m *MockInterface) Tx() (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx")
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Tx indicates an expected call of Tx.
func (mr *MockInterfaceMockRecorder) Tx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockInterface)(nil).Tx))
}
