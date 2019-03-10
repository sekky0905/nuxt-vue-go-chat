package application

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
)

// AuthenticationService is the interface of AuthenticationService.
type AuthenticationService interface {
	SignUp(ctx context.Context, name, password string) (*model.User, error)
}

// authenticationService is the service of authentication.
type authenticationService struct {
	m        repository.SQLManager
	uFactory service.UserRepoFactory
	sFactory service.SessionRepoFactory
	txCloser CloseTransaction
}

// NewAuthenticationApplication generates and returns AuthenticationService.
func NewAuthenticationApplication(m repository.SQLManager, uFactory service.UserRepoFactory, sFactory service.SessionRepoFactory, txCloser CloseTransaction) AuthenticationService {
	return &authenticationService{
		m:        m,
		uFactory: uFactory,
		sFactory: sFactory,
		txCloser: txCloser,
	}
}

// SignUp sign up an user.
func (a *authenticationService) SignUp(ctx context.Context, name, password string) (user *model.User, err error) {
	return nil, nil
}
