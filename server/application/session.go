package application

import (
	"github.com/sekky0905/go-vue-chat/server/domain/repository"
)

// SessionService is the interface of SessionService.
type SessionService interface {
}

// sessionService is the service of user.
type sessionService struct {
	m        repository.SQLManager
	sFactory repository.SessionRepoFactory
	txCloser CloseTransaction
}

// NewSessionService generates and returns NewSessionService.
func NewSessionService(m repository.SQLManager, f repository.SessionRepoFactory, txCloser CloseTransaction) SessionService {
	return &sessionService{
		m:        m,
		sFactory: f,
		txCloser: txCloser,
	}
}
