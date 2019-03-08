package db

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	log "github.com/sirupsen/logrus"
)

// sessionRepository is repository of user.
type sessionRepository struct {
	ctx context.Context
}

// NewSessionRepository generates and returns sessionRepository.
func NewSessionRepository(ctx context.Context) repository.SessionRepository {
	return &sessionRepository{
		ctx: ctx,
	}
}

// ErrorMsg generates and returns error message.
func (repo *sessionRepository) ErrorMsg(method model.RepositoryMethod, err error) error {
	log.Error(err.Error())
	return &model.RepositoryError{
		BaseErr:                     err,
		RepositoryMethod:            method,
		DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
		DomainModelNameForUser:      model.DomainModelNameSessionForUser,
	}
}

// GetSessionByID gets and returns a record specified by id.
func (repo *sessionRepository) GetSessionByID(m repository.DBManager, id string) (*model.Session, error) {
	return nil, nil
}

// list gets and returns list of records.
func (repo *sessionRepository) list(m repository.DBManager, method model.RepositoryMethod, query string, args ...interface{}) (sessions []*model.Session, err error) {
	return nil, nil
}

// InsertSession insert a record.
func (repo *sessionRepository) InsertSession(m repository.DBManager, session *model.Session) (err error) {
	return nil
}

// DeleteSession delete a record.
func (repo *sessionRepository) DeleteSession(m repository.DBManager, id string) error {

	return nil
}
