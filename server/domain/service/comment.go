package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// ThreadService is interface of ThreadService.
type CommentService interface {
	IsAlreadyExistID(ctx context.Context, id uint32) (bool, error)
}

// threadService is domain service of Thread.
type commentService struct {
	repo repository.CommentRepository
	m    repository.SQLManager
}

// NewThreadService generates and returns ThreadService.
func NewCommentService(repo repository.CommentRepository, m repository.SQLManager) CommentService {
	return &commentService{
		repo: repo,
		m:    m,
	}
}

// NewComment generates and returns ThreadService.
func NewComment(content string, threadID uint32, user *model.User) *model.Comment {
	return &model.Comment{
		Content:  content,
		ThreadID: threadID,
		User:     user,
	}
}

// NewCommentList generates and returns ThreadService.
func NewCommentList(list []*model.Comment, hasNext bool, cursor uint32) *model.CommentList {
	return &model.CommentList{
		Comments: list,
		HasNext:  hasNext,
		Cursor:   cursor,
	}
}

// IsAlreadyExistID checks duplication of id.
func (s commentService) IsAlreadyExistID(ctx context.Context, id uint32) (bool, error) {
	searched, err := s.repo.GetCommentByID(ctx, s.m, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to get comment by PropertyNameForDeveloper")
	}
	return searched != nil, nil
}
