// Code generated by MockGen. DO NOT EDIT.
// Source: internal/gophermart/repository/user_withdraw.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	models "github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

// MockRepositoryCategoryUserWithdraw is a mock of RepositoryCategoryUserWithdraw interface.
type MockRepositoryCategoryUserWithdraw struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryCategoryUserWithdrawMockRecorder
}

// MockRepositoryCategoryUserWithdrawMockRecorder is the mock recorder for MockRepositoryCategoryUserWithdraw.
type MockRepositoryCategoryUserWithdrawMockRecorder struct {
	mock *MockRepositoryCategoryUserWithdraw
}

// NewMockRepositoryCategoryUserWithdraw creates a new mock instance.
func NewMockRepositoryCategoryUserWithdraw(ctrl *gomock.Controller) *MockRepositoryCategoryUserWithdraw {
	mock := &MockRepositoryCategoryUserWithdraw{ctrl: ctrl}
	mock.recorder = &MockRepositoryCategoryUserWithdrawMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryCategoryUserWithdraw) EXPECT() *MockRepositoryCategoryUserWithdrawMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepositoryCategoryUserWithdraw) Create(ctx context.Context, withdraw *models.UserWithdraw) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, withdraw)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryCategoryUserWithdrawMockRecorder) Create(ctx, withdraw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepositoryCategoryUserWithdraw)(nil).Create), ctx, withdraw)
}

// GetAllByUserID mocks base method.
func (m *MockRepositoryCategoryUserWithdraw) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserWithdraw, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUserID", ctx, userID)
	ret0, _ := ret[0].([]models.UserWithdraw)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByUserID indicates an expected call of GetAllByUserID.
func (mr *MockRepositoryCategoryUserWithdrawMockRecorder) GetAllByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUserID", reflect.TypeOf((*MockRepositoryCategoryUserWithdraw)(nil).GetAllByUserID), ctx, userID)
}

// GetWithdrawnSumByUserID mocks base method.
func (m *MockRepositoryCategoryUserWithdraw) GetWithdrawnSumByUserID(ctx context.Context, userID uuid.UUID) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawnSumByUserID", ctx, userID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawnSumByUserID indicates an expected call of GetWithdrawnSumByUserID.
func (mr *MockRepositoryCategoryUserWithdrawMockRecorder) GetWithdrawnSumByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawnSumByUserID", reflect.TypeOf((*MockRepositoryCategoryUserWithdraw)(nil).GetWithdrawnSumByUserID), ctx, userID)
}