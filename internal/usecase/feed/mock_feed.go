// Code generated by MockGen. DO NOT EDIT.
// Source: feed.go
//
// Generated by this command:
//
//	mockgen -source=feed.go -destination=./mock_feed.go -package=feed
//

// Package feed is a generated GoMock package.
package feed

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	entity "github.com/soltanat/otus-highload/internal/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockPostStorager is a mock of PostStorager interface.
type MockPostStorager struct {
	ctrl     *gomock.Controller
	recorder *MockPostStoragerMockRecorder
}

// MockPostStoragerMockRecorder is the mock recorder for MockPostStorager.
type MockPostStoragerMockRecorder struct {
	mock *MockPostStorager
}

// NewMockPostStorager creates a new mock instance.
func NewMockPostStorager(ctrl *gomock.Controller) *MockPostStorager {
	mock := &MockPostStorager{ctrl: ctrl}
	mock.recorder = &MockPostStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostStorager) EXPECT() *MockPostStoragerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockPostStorager) List(ctx context.Context, tx entity.Tx, filter *entity.PostFilter) ([]entity.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, tx, filter)
	ret0, _ := ret[0].([]entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockPostStoragerMockRecorder) List(ctx, tx, filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPostStorager)(nil).List), ctx, tx, filter)
}

// MockFriendStorager is a mock of FriendStorager interface.
type MockFriendStorager struct {
	ctrl     *gomock.Controller
	recorder *MockFriendStoragerMockRecorder
}

// MockFriendStoragerMockRecorder is the mock recorder for MockFriendStorager.
type MockFriendStoragerMockRecorder struct {
	mock *MockFriendStorager
}

// NewMockFriendStorager creates a new mock instance.
func NewMockFriendStorager(ctrl *gomock.Controller) *MockFriendStorager {
	mock := &MockFriendStorager{ctrl: ctrl}
	mock.recorder = &MockFriendStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFriendStorager) EXPECT() *MockFriendStoragerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockFriendStorager) List(ctx context.Context, tx entity.Tx, filter *entity.FriendFilter) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, tx, filter)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockFriendStoragerMockRecorder) List(ctx, tx, filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockFriendStorager)(nil).List), ctx, tx, filter)
}

// MockUserStorager is a mock of UserStorager interface.
type MockUserStorager struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoragerMockRecorder
}

// MockUserStoragerMockRecorder is the mock recorder for MockUserStorager.
type MockUserStoragerMockRecorder struct {
	mock *MockUserStorager
}

// NewMockUserStorager creates a new mock instance.
func NewMockUserStorager(ctrl *gomock.Controller) *MockUserStorager {
	mock := &MockUserStorager{ctrl: ctrl}
	mock.recorder = &MockUserStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStorager) EXPECT() *MockUserStoragerMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockUserStorager) Find(ctx context.Context, tx entity.Tx, filter *entity.UserFilter) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, tx, filter)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockUserStoragerMockRecorder) Find(ctx, tx, filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUserStorager)(nil).Find), ctx, tx, filter)
}

// Get mocks base method.
func (m *MockUserStorager) Get(ctx context.Context, tx entity.Tx, userID uuid.UUID) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, tx, userID)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserStoragerMockRecorder) Get(ctx, tx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserStorager)(nil).Get), ctx, tx, userID)
}
