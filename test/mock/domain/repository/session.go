// Code generated by MockGen. DO NOT EDIT.
// Source: session.go
//
// Generated by this command:
//
//	mockgen -source=session.go -package=repository -destination=../../../../../test/mock/domain/repository/session.go
//

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	entity "github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockSessionRepository is a mock of SessionRepository interface.
type MockSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepositoryMockRecorder
	isgomock struct{}
}

// MockSessionRepositoryMockRecorder is the mock recorder for MockSessionRepository.
type MockSessionRepositoryMockRecorder struct {
	mock *MockSessionRepository
}

// NewMockSessionRepository creates a new mock instance.
func NewMockSessionRepository(ctrl *gomock.Controller) *MockSessionRepository {
	mock := &MockSessionRepository{ctrl: ctrl}
	mock.recorder = &MockSessionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRepository) EXPECT() *MockSessionRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockSessionRepository) Delete(arg0 context.Context, arg1 *entity.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSessionRepositoryMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessionRepository)(nil).Delete), arg0, arg1)
}

// FindOneByAccountID mocks base method.
func (m *MockSessionRepository) FindOneByAccountID(arg0 context.Context, arg1 uuid.UUID) (*entity.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByAccountID", arg0, arg1)
	ret0, _ := ret[0].(*entity.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByAccountID indicates an expected call of FindOneByAccountID.
func (mr *MockSessionRepositoryMockRecorder) FindOneByAccountID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByAccountID", reflect.TypeOf((*MockSessionRepository)(nil).FindOneByAccountID), arg0, arg1)
}

// FindOneByTokenAndNotExpired mocks base method.
func (m *MockSessionRepository) FindOneByTokenAndNotExpired(arg0 context.Context, arg1 string) (*entity.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByTokenAndNotExpired", arg0, arg1)
	ret0, _ := ret[0].(*entity.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByTokenAndNotExpired indicates an expected call of FindOneByTokenAndNotExpired.
func (mr *MockSessionRepositoryMockRecorder) FindOneByTokenAndNotExpired(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByTokenAndNotExpired", reflect.TypeOf((*MockSessionRepository)(nil).FindOneByTokenAndNotExpired), arg0, arg1)
}

// Save mocks base method.
func (m *MockSessionRepository) Save(arg0 context.Context, arg1 *entity.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockSessionRepositoryMockRecorder) Save(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSessionRepository)(nil).Save), arg0, arg1)
}
