package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/util"
)

// AuthenticationService is interface of domain service of authentication.
type AuthenticationService interface {
	Authenticate(ctx context.Context, userName, password string) (ok bool, user *model.User, err error)
}

// authenticationService is domain service of authentication.
type authenticationService struct {
	m    repository.DBManager
	repo repository.UserRepository
}

// NewAuthenticationService generates and returns AuthenticationService.
func NewAuthenticationService(m repository.DBManager, repo repository.UserRepository) AuthenticationService {
	return &authenticationService{
		m:    m,
		repo: repo,
	}
}

// Authenticate authenticate user.
func (s *authenticationService) Authenticate(ctx context.Context, userName, password string) (ok bool, user *model.User, err error) {
	gotUser, err := s.repo.GetUserByName(ctx, s.m, userName)
	if err != nil {
		if _, ok := errors.Cause(err).(*model.NoSuchDataError); ok {
			return false, nil, &model.AuthenticationErr{
				BaseErr: err,
			}
		}

		return false, nil, errors.Wrap(err, "failed, to get user by name")
	}

	if !util.CheckHashOfPassword(password, gotUser.Password) {
		return false, nil, nil
	}

	return true, gotUser, nil
}
