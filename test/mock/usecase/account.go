// Code generated by MockGen. DO NOT EDIT.
// Source: account.go
//
// Generated by this command:
//
//	mockgen -source=account.go -package=usecase -destination=../../../../test/mock/usecase/account.go
//

// Package usecase is a generated GoMock package.
package usecase

import (
	context "context"
	reflect "reflect"

	dto "github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockAccountUsecase is a mock of AccountUsecase interface.
type MockAccountUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAccountUsecaseMockRecorder
	isgomock struct{}
}

// MockAccountUsecaseMockRecorder is the mock recorder for MockAccountUsecase.
type MockAccountUsecaseMockRecorder struct {
	mock *MockAccountUsecase
}

// NewMockAccountUsecase creates a new mock instance.
func NewMockAccountUsecase(ctrl *gomock.Controller) *MockAccountUsecase {
	mock := &MockAccountUsecase{ctrl: ctrl}
	mock.recorder = &MockAccountUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountUsecase) EXPECT() *MockAccountUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAccountUsecase) Create(arg0 context.Context, arg1, arg2, arg3 string) (*dto.AccountDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*dto.AccountDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAccountUsecaseMockRecorder) Create(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAccountUsecase)(nil).Create), arg0, arg1, arg2, arg3)
}
