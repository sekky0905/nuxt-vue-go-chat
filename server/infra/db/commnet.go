package db

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	. "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// commentRepository is repository of comment.
type commentRepository struct {
}

// NewCommentRepository generates and returns CommentRepository.
func NewCommentRepository() CommentRepository {
	return &commentRepository{}
}

// ErrorMsg generates and returns error message.
func (repo *commentRepository) ErrorMsg(method model.RepositoryMethod, err error) error {
	return &model.RepositoryError{
		BaseErr:                     err,
		RepositoryMethod:            method,
		DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
		DomainModelNameForUser:      model.DomainModelNameCommentForUser,
	}
}

// ListThreads lists ThreadList.
func (repo *commentRepository) ListComments(ctx context.Context, m DBManager, threadID uint32, limit int, cursor uint32) (*model.CommentList, error) {
	return nil, nil
}

// GetThreadByID gets and returns a record specified by id.
func (repo *commentRepository) GetCommentByID(ctx context.Context, m DBManager, id uint32) (*model.Comment, error) {
	return nil, nil
}

// InsertThread insert a record.
func (repo *commentRepository) InsertComment(ctx context.Context, m DBManager, user *model.Comment) (uint32, error) {
	return 1, nil
}

// UpdateComment updates a record.
func (repo *commentRepository) UpdateComment(ctx context.Context, m DBManager, id uint32, thead *model.Comment) error {
	return nil
}

// DeleteComment delete a record.
func (repo *commentRepository) DeleteComment(ctx context.Context, m DBManager, id uint32) error {
	return nil
}
