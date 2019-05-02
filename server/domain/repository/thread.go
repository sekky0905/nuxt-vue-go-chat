package repository

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"
)

// ThreadRepository is Repository of Thread.
type ThreadRepository interface {
	ListThreads(ctx context.Context, m query.SQLManager, cursor uint32, limit int) (*model.ThreadList, error)
	GetThreadByID(ctx context.Context, m query.SQLManager, id uint32) (*model.Thread, error)
	GetThreadByTitle(ctx context.Context, m query.SQLManager, name string) (*model.Thread, error)
	InsertThread(ctx context.Context, m query.SQLManager, thead *model.Thread) (uint32, error)
	UpdateThread(ctx context.Context, m query.SQLManager, id uint32, thead *model.Thread) error
	DeleteThread(ctx context.Context, m query.SQLManager, id uint32) error
}
