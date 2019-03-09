package db

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

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
	ex := &model.RepositoryError{
		BaseErr:                     err,
		RepositoryMethod:            method,
		DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
		DomainModelNameForUser:      model.DomainModelNameSessionForUser,
	}

	log.Errorf("EX--->%#v", ex)
	return ex
}

// GetSessionByID gets and returns a record specified by id.
func (repo *sessionRepository) GetSessionByID(m repository.DBManager, id string) (*model.Session, error) {
	query := "SELECT id, user_id, created_at, updated_at FROM sessions WHERE id=?"

	list, err := repo.list(m, model.RepositoryMethodREAD, query, id)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:                     err,
			PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
			PropertyNameForUser:         model.IDPropertyForUser,
			PropertyValue:               id,
			DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameSessionForUser,
		}
		return nil, errors.WithStack(err)
	}

	if err != nil {
		return nil, repo.ErrorMsg(model.RepositoryMethodREAD, errors.WithStack(err))
	}

	return list[0], nil
}

// list gets and returns list of records.
func (repo *sessionRepository) list(m repository.DBManager, method model.RepositoryMethod, query string, args ...interface{}) (sessions []*model.Session, err error) {
	stmt, err := m.PrepareContext(repo.ctx, query)
	if err != nil {
		return nil, repo.ErrorMsg(method, errors.WithStack(err))
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	rows, err := stmt.QueryContext(repo.ctx, args...)
	if err != nil {
		err = repo.ErrorMsg(method, errors.WithStack(err))
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	list := make([]*model.Session, 0)
	for rows.Next() {
		session := &model.Session{}

		err = rows.Scan(
			&session.ID,
			&session.UserID,
			&session.CreatedAt,
			&session.UpdatedAt,
		)

		if err != nil {
			return nil, repo.ErrorMsg(method, errors.WithStack(err))
		}

		list = append(list, session)
	}

	return list, nil
}

// InsertSession insert a record.
func (repo *sessionRepository) InsertSession(m repository.DBManager, session *model.Session) error {
	query := "INSERT INTO sessions (id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)"
	stmt, err := m.PrepareContext(repo.ctx, query)
	if err != nil {
		return errors.WithStack(repo.ErrorMsg(model.RepositoryMethodInsert, err))
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	result, err := stmt.ExecContext(repo.ctx, session.ID, session.UserID, session.CreatedAt, session.UpdatedAt)
	if err != nil {
		log.Errorf("[]==%s", err.Error())
		return errors.WithStack(repo.ErrorMsg(model.RepositoryMethodInsert, err))
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = fmt.Errorf("total affected: %d ", affect)
		return errors.WithStack(repo.ErrorMsg(model.RepositoryMethodInsert, err))
	}

	return nil
}

// DeleteSession delete a record.
func (repo *sessionRepository) DeleteSession(m repository.DBManager, id string) error {

	return nil
}
