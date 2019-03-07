package db

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	. "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// userRepository is repository of user.
type userRepository struct {
	ctx context.Context
}

// NewUserRepository は、userRepository生成する。
func NewUserRepository(ctx context.Context) UserRepository {
	return &userRepository{
		ctx: ctx,
	}
}

// ErrorMsg generates and returns error message.
func (repo *userRepository) ErrorMsg(method RepositoryMethod, err error) error {
	return nil
}

func (repo *userRepository) GetUserByID(m DBManager, id uint32) (*model.User, error)     {}
func (repo *userRepository) GetUserByName(m DBManager, name string) (*model.User, error) {}
func (repo *userRepository) InsertUser(m DBManager, user *model.User) (uint32, error)    {}
func (repo *userRepository) UpdateUser(m DBManager, id uint32, user *model.User) error   {}
func (repo *userRepository) DeleteUser(m DBManager, id uint32) error                     {}
