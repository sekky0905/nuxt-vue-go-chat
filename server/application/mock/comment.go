// Code generated by MockGen. DO NOT EDIT.
// Source: server/application/comment.go

// Package mock_application is a generated GoMock package.
package mock_application

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	reflect "reflect"
)

// MockCommentService is a mock of CommentService interface
type MockCommentService struct {
	ctrl     *gomock.Controller
	recorder *MockCommentServiceMockRecorder
}

// MockCommentServiceMockRecorder is the mock recorder for MockCommentService
type MockCommentServiceMockRecorder struct {
	mock *MockCommentService
}

// NewMockCommentService creates a new mock instance
func NewMockCommentService(ctrl *gomock.Controller) *MockCommentService {
	mock := &MockCommentService{ctrl: ctrl}
	mock.recorder = &MockCommentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommentService) EXPECT() *MockCommentServiceMockRecorder {
	return m.recorder
}

// ListComments mocks base method
func (m *MockCommentService) ListComments(ctx context.Context, threadID uint32, limit int, cursor uint32) (*model.CommentList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListComments", ctx, threadID, limit, cursor)
	ret0, _ := ret[0].(*model.CommentList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListComments indicates an expected call of ListComments
func (mr *MockCommentServiceMockRecorder) ListComments(ctx, threadID, limit, cursor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListComments", reflect.TypeOf((*MockCommentService)(nil).ListComments), ctx, threadID, limit, cursor)
}

// GetComment mocks base method
func (m *MockCommentService) GetComment(ctx context.Context, id uint32) (*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComment", ctx, id)
	ret0, _ := ret[0].(*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetComment indicates an expected call of GetComment
func (mr *MockCommentServiceMockRecorder) GetComment(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComment", reflect.TypeOf((*MockCommentService)(nil).GetComment), ctx, id)
}

// CreateComment mocks base method
func (m *MockCommentService) CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateComment", ctx, comment)
	ret0, _ := ret[0].(*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateComment indicates an expected call of CreateComment
func (mr *MockCommentServiceMockRecorder) CreateComment(ctx, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateComment", reflect.TypeOf((*MockCommentService)(nil).CreateComment), ctx, comment)
}

// UpdateComment mocks base method
func (m *MockCommentService) UpdateComment(ctx context.Context, id uint32, comment *model.Comment) (*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComment", ctx, id, comment)
	ret0, _ := ret[0].(*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateComment indicates an expected call of UpdateComment
func (mr *MockCommentServiceMockRecorder) UpdateComment(ctx, id, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComment", reflect.TypeOf((*MockCommentService)(nil).UpdateComment), ctx, id, comment)
}

// DeleteComment mocks base method
func (m *MockCommentService) DeleteComment(ctx context.Context, id uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteComment", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteComment indicates an expected call of DeleteComment
func (mr *MockCommentServiceMockRecorder) DeleteComment(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComment", reflect.TypeOf((*MockCommentService)(nil).DeleteComment), ctx, id)
}
