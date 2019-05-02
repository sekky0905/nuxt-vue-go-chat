package repository

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"
)

// UserRepository is repository of user.
type UserRepository interface {
	GetUserByID(ctx context.Context, m query.SQLManager, id uint32) (*model.User, error)
	GetUserByName(ctx context.Context, m query.SQLManager, name string) (*model.User, error)
	InsertUser(ctx context.Context, m query.SQLManager, user *model.User) (uint32, error)
	UpdateUser(ctx context.Context, m query.SQLManager, id uint32, user *model.User) error
	DeleteUser(ctx context.Context, m query.SQLManager, id uint32) error
}
