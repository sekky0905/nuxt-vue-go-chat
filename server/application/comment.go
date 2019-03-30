package application

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	. "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
)

// ThreadService is interface of ThreadService.
type CommentService interface {
	ListComments(ctx context.Context, threadId uint32, limit int, cursor uint32) (*model.CommentList, error)
	GetComment(ctx context.Context, id uint32) (*model.Comment, error)
	CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	UpdateComment(ctx context.Context, id uint32, comment *model.Comment) (*model.Comment, error)
	DeleteComment(ctx context.Context, id uint32) error
}

// threadService is application service of thread.
type commentApplication struct {
	m        DBManager
	service  service.CommentService
	repo     CommentRepository
	txCloser CloseTransaction
}

// NewThreadService generates and returns ThreadService.
func NewCommentApplication(m DBManager, service service.CommentService, repo CommentRepository, txCloser CloseTransaction) CommentService {
	return &commentApplication{
		m:        m,
		service:  service,
		repo:     repo,
		txCloser: txCloser,
	}
}

// ListThreads gets ThreadList.
func (a *commentApplication) ListComments(ctx context.Context, threadId uint32, limit int, cursor uint32) (*model.CommentList, error) {

	return nil, nil
}

// GetThread gets Thread.
func (a *commentApplication) GetComment(ctx context.Context, id uint32) (*model.Comment, error) {

	return nil, nil
}

// CreateThread creates Thread.
func (a *commentApplication) CreateComment(ctx context.Context, param *model.Comment) (comment *model.Comment, err error) {

	return nil, nil
}

// UpdateThread updates Thread.
func (a *commentApplication) UpdateComment(ctx context.Context, id uint32, param *model.Comment) (comment *model.Comment, err error) {
	return param, nil
}

// DeleteThread deletes Thread.
func (a *commentApplication) DeleteComment(ctx context.Context, id uint32) (err error) {

	return nil
}
