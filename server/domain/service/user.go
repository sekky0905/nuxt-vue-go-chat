package service

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// UserService is interface of domain service of user.
type UserService interface {
	IsAlreadyExistID(ctx context.Context, id uint32) (bool, error)
	IsAlreadyExistName(ctx context.Context, name string) (bool, error)
}

// userService is domain service of user.
type userService struct {
	repo repository.UserRepository
	m    repository.SQLManager
}

// NewUserService generates and returns UserService.
func NewUserService(repo repository.UserRepository, m repository.SQLManager) UserService {
	return &userService{
		repo: repo,
		m:    m,
	}
}

// IsAlreadyExistID checks whether the data specified by id already exists or not.
func (s *userService) IsAlreadyExistID(ctx context.Context, id uint32) (bool, error) {
	panic("implement me")
}

// IsAlreadyExistName checks whether the data specified by name already exists or not.
func (s *userService) IsAlreadyExistName(ctx context.Context, name string) (bool, error) {
	panic("implement me")
}
