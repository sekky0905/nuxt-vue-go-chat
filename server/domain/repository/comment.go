package repository

import "github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"

// CommentRepository is Repository of Comment.
type CommentRepository interface {
	ListComments(m DBManager, threadId uint32, limit int, cursor uint32) (*model.CommentList, error)
	GetCommentByID(m DBManager, id uint32) (*model.Comment, error)
	InsertComment(m DBManager, user *model.Comment) (uint32, error)
	UpdateComment(m DBManager, id uint32, thead *model.Comment) error
	DeleteComment(m DBManager, id uint32) error
}
