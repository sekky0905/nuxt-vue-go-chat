package application

import (
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
)

// UserService is the interface of UserService.
type UserService interface {
}

// userService is the service of user.
type userService struct {
	m        repository.DBManager
	uFactory service.UserRepoFactory
	txCloser CloseTransaction
}

// NewUserService generates and returns UserService.
func NewUserService(m repository.DBManager, f service.UserRepoFactory, txCloser CloseTransaction) UserService {
	return &userService{
		m:        m,
		uFactory: f,
		txCloser: txCloser,
	}
}
