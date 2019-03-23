package repository

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// ThreadRepository is Repository of Thread.
type ThreadRepository interface {
	ListThreads(ctx context.Context, m SQLManager, cursor uint32, limit int) (*model.ThreadList, error)
	GetThreadByID(ctx context.Context, m SQLManager, id uint32) (*model.Thread, error)
	GetThreadByTitle(ctx context.Context, m SQLManager, name string) (*model.Thread, error)
	InsertThread(ctx context.Context, m SQLManager, user *model.Thread) (uint32, error)
	UpdateThread(ctx context.Context, m SQLManager, id uint32, thead *model.Thread) error
	DeleteThread(ctx context.Context, m SQLManager, id uint32) error
}

// ThreadRepoFactory は、ThreadRepositoryのFactory。
// ThreadRepository is Factory of ThreadRepository.
type ThreadRepoFactory func(ctx context.Context) ThreadRepository
