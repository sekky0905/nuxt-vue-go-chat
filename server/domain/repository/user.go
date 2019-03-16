package repository

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// UserRepository is repository of user.
type UserRepository interface {
	GetUserByID(m SQLManager, id uint32) (*model.User, error)
	GetUserByName(m SQLManager, name string) (*model.User, error)
	InsertUser(m SQLManager, user *model.User) (uint32, error)
	UpdateUser(m SQLManager, id uint32, user *model.User) error
	DeleteUser(m SQLManager, id uint32) error
}

// UserRepoFactory is Factory of UserRepository.
type UserRepoFactory func(ctx context.Context) UserRepository
