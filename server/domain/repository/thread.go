package repository

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// ThreadRepository is Repository of Thread.
type ThreadRepository interface {
	ListThreads(m DBManager, limit int, cursor uint32) (*model.ThreadList, error)
	GetThreadByID(m DBManager, id uint32) (*model.Thread, error)
	GetThreadByTitle(m DBManager, name string) (*model.Thread, error)
	InsertThread(m DBManager, user *model.Thread) (uint32, error)
	UpdateThread(m DBManager, id uint32, thead *model.Thread) error
	DeleteThread(m DBManager, id uint32) error
}

// ThreadRepoFactory は、ThreadRepositoryのFactory。
// ThreadRepository is Factory of ThreadRepository.
type ThreadRepoFactory func(ctx context.Context) ThreadRepository
