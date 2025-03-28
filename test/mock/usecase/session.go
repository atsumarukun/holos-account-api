// Code generated by MockGen. DO NOT EDIT.
// Source: session.go
//
// Generated by this command:
//
//	mockgen -source=session.go -package=usecase -destination=../../../../test/mock/usecase/session.go
//

// Package usecase is a generated GoMock package.
package usecase

import (
	context "context"
	reflect "reflect"

	dto "github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockSessionUsecase is a mock of SessionUsecase interface.
type MockSessionUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockSessionUsecaseMockRecorder
	isgomock struct{}
}

// MockSessionUsecaseMockRecorder is the mock recorder for MockSessionUsecase.
type MockSessionUsecaseMockRecorder struct {
	mock *MockSessionUsecase
}

// NewMockSessionUsecase creates a new mock instance.
func NewMockSessionUsecase(ctrl *gomock.Controller) *MockSessionUsecase {
	mock := &MockSessionUsecase{ctrl: ctrl}
	mock.recorder = &MockSessionUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionUsecase) EXPECT() *MockSessionUsecaseMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockSessionUsecase) Authenticate(arg0 context.Context, arg1 string) (*dto.AccountDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", arg0, arg1)
	ret0, _ := ret[0].(*dto.AccountDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockSessionUsecaseMockRecorder) Authenticate(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockSessionUsecase)(nil).Authenticate), arg0, arg1)
}

// Authorize mocks base method.
func (m *MockSessionUsecase) Authorize(arg0 context.Context, arg1 uuid.UUID) (*dto.AccountDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", arg0, arg1)
	ret0, _ := ret[0].(*dto.AccountDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorize indicates an expected call of Authorize.
func (mr *MockSessionUsecaseMockRecorder) Authorize(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockSessionUsecase)(nil).Authorize), arg0, arg1)
}

// Login mocks base method.
func (m *MockSessionUsecase) Login(arg0 context.Context, arg1, arg2 string) (*dto.SessionDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1, arg2)
	ret0, _ := ret[0].(*dto.SessionDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockSessionUsecaseMockRecorder) Login(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockSessionUsecase)(nil).Login), arg0, arg1, arg2)
}

// Logout mocks base method.
func (m *MockSessionUsecase) Logout(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockSessionUsecaseMockRecorder) Logout(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockSessionUsecase)(nil).Logout), arg0, arg1)
}
