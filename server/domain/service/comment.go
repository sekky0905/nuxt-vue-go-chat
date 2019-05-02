package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"
)

// CommentService is interface of CommentService.
type CommentService interface {
	IsAlreadyExistID(ctx context.Context, m query.SQLManager, id uint32) (bool, error)
}

// commentService is domain service of Comment.
type commentService struct {
	repo repository.CommentRepository
}

// NewCommentService generates and returns CommentService.
func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}

// NewComment generates and returns CommentService.
func NewComment(content string, threadID uint32, user *model.User) *model.Comment {
	return &model.Comment{
		Content:  content,
		ThreadID: threadID,
		User:     user,
	}
}

// NewCommentList generates and returns CommentService.
func NewCommentList(list []*model.Comment, hasNext bool, cursor uint32) *model.CommentList {
	return &model.CommentList{
		Comments: list,
		HasNext:  hasNext,
		Cursor:   cursor,
	}
}

// IsAlreadyExistID checks duplication of id.
func (s commentService) IsAlreadyExistID(ctx context.Context, m query.SQLManager, id uint32) (bool, error) {
	searched, err := s.repo.GetCommentByID(ctx, m, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to get comment by PropertyNameForDeveloper")
	}
	return searched != nil, nil
}
