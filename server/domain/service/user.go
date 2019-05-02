package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"
	"github.com/sekky0905/nuxt-vue-go-chat/server/util"
)

// UserService is interface of domain service of user.
type UserService interface {
	NewUser(name, password string) (*model.User, error)
	IsAlreadyExistID(ctx context.Context, m query.SQLManager, id uint32) (bool, error)
	IsAlreadyExistName(ctx context.Context, m query.SQLManager, name string) (bool, error)
}

// UserRepoFactory is factory of UserRepository.
type UserRepoFactory func(ctx context.Context) repository.UserRepository

// userService is domain service of user.
type userService struct {
	repo repository.UserRepository
}

// NewUserService generates and returns UserService.
func NewUserService(m query.SQLManager, repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// NewUser generates and reruns User.
func (s *userService) NewUser(name, password string) (*model.User, error) {
	hashed, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Name:      name,
		Password:  hashed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// IsAlreadyExistID checks whether the data specified by id already exists or not.
func (s *userService) IsAlreadyExistID(ctx context.Context, m query.SQLManager, id uint32) (bool, error) {
	searched, err := s.repo.GetUserByID(ctx, m, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user by id")
	}
	return searched != nil, nil
}

// IsAlreadyExistName checks whether the data specified by name already exists or not.
func (s *userService) IsAlreadyExistName(ctx context.Context, m query.SQLManager, name string) (bool, error) {
	searched, err := s.repo.GetUserByName(ctx, m, name)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user by Name")
	}

	return searched != nil, nil
}
