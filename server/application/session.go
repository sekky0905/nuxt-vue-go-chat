package application

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
	"github.com/sekky0905/nuxt-vue-go-chat/server/util"
)

// SessionService is the interface of SessionService.
type SessionService interface {
}

// sessionService is the service of user.
type sessionService struct {
	m        repository.SQLManager
	sFactory service.SessionRepoFactory
	txCloser CloseTransaction
}

// NewSessionService generates and returns NewSessionService.
func NewSessionService(m repository.SQLManager, f service.SessionRepoFactory, txCloser CloseTransaction) SessionService {
	return &sessionService{
		m:        m,
		sFactory: f,
		txCloser: txCloser,
	}
}

// createSession creates the session.
func createSession(ctx context.Context, userID uint32, db repository.DBManager, repo repository.SessionRepository, sService service.SessionService) (*model.Session, error) {
	session := model.NewSession(userID)
	session.ID = util.UUID()

	// ready for collision of UUID.
	yes := true
	var err error
	for yes {
		yes, err = sService.IsAlreadyExistID(ctx, session.ID)
		if err != nil {
			if _, ok := errors.Cause(err).(*model.NoSuchDataError); !ok {
				return nil, errors.Wrap(err, "failed to check whether already exists id or not")
			}
		}
	}

	if err := repo.InsertSession(db, session); err != nil {
		return nil, errors.Wrap(err, "failed to insert session")
	}
	return session, nil
}
