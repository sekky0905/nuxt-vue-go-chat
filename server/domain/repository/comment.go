package repository

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// CommentRepository is Repository of Comment.
type CommentRepository interface {
	ListComments(ctx context.Context, m SQLManager, threadID uint32, limit int, cursor uint32) (*model.CommentList, error)
	GetCommentByID(ctx context.Context, m SQLManager, id uint32) (*model.Comment, error)
	InsertComment(ctx context.Context, m SQLManager, user *model.Comment) (uint32, error)
	UpdateComment(ctx context.Context, m SQLManager, id uint32, thead *model.Comment) error
	DeleteComment(ctx context.Context, m SQLManager, id uint32) error
}
