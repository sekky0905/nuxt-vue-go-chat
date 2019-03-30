// Code generated by MockGen. DO NOT EDIT.
// Source: domain/repository/comment.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	repository "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	reflect "reflect"
)

// MockCommentRepository is a mock of CommentRepository interface
type MockCommentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCommentRepositoryMockRecorder
}

// MockCommentRepositoryMockRecorder is the mock recorder for MockCommentRepository
type MockCommentRepositoryMockRecorder struct {
	mock *MockCommentRepository
}

// NewMockCommentRepository creates a new mock instance
func NewMockCommentRepository(ctrl *gomock.Controller) *MockCommentRepository {
	mock := &MockCommentRepository{ctrl: ctrl}
	mock.recorder = &MockCommentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommentRepository) EXPECT() *MockCommentRepositoryMockRecorder {
	return m.recorder
}

// ListComments mocks base method
func (m_2 *MockCommentRepository) ListComments(ctx context.Context, m repository.SQLManager, threadID uint32, limit int, cursor uint32) (*model.CommentList, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "ListComments", ctx, m, threadID, limit, cursor)
	ret0, _ := ret[0].(*model.CommentList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListComments indicates an expected call of ListComments
func (mr *MockCommentRepositoryMockRecorder) ListComments(ctx, m, threadID, limit, cursor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListComments", reflect.TypeOf((*MockCommentRepository)(nil).ListComments), ctx, m, threadID, limit, cursor)
}

// GetCommentByID mocks base method
func (m_2 *MockCommentRepository) GetCommentByID(ctx context.Context, m repository.SQLManager, id uint32) (*model.Comment, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "GetCommentByID", ctx, m, id)
	ret0, _ := ret[0].(*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentByID indicates an expected call of GetCommentByID
func (mr *MockCommentRepositoryMockRecorder) GetCommentByID(ctx, m, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentByID", reflect.TypeOf((*MockCommentRepository)(nil).GetCommentByID), ctx, m, id)
}

// InsertComment mocks base method
func (m_2 *MockCommentRepository) InsertComment(ctx context.Context, m repository.SQLManager, comment *model.Comment) (uint32, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "InsertComment", ctx, m, comment)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertComment indicates an expected call of InsertComment
func (mr *MockCommentRepositoryMockRecorder) InsertComment(ctx, m, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertComment", reflect.TypeOf((*MockCommentRepository)(nil).InsertComment), ctx, m, comment)
}

// UpdateComment mocks base method
func (m_2 *MockCommentRepository) UpdateComment(ctx context.Context, m repository.SQLManager, id uint32, comment *model.Comment) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "UpdateComment", ctx, m, id, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateComment indicates an expected call of UpdateComment
func (mr *MockCommentRepositoryMockRecorder) UpdateComment(ctx, m, id, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComment", reflect.TypeOf((*MockCommentRepository)(nil).UpdateComment), ctx, m, id, comment)
}

// DeleteComment mocks base method
func (m_2 *MockCommentRepository) DeleteComment(ctx context.Context, m repository.SQLManager, id uint32) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "DeleteComment", ctx, m, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteComment indicates an expected call of DeleteComment
func (mr *MockCommentRepositoryMockRecorder) DeleteComment(ctx, m, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComment", reflect.TypeOf((*MockCommentRepository)(nil).DeleteComment), ctx, m, id)
}
