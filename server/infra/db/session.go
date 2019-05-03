package db

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"
	"go.uber.org/zap"
)

// sessionRepository is repository of user.
type sessionRepository struct {
}

// NewSessionRepository generates and returns sessionRepository.
func NewSessionRepository() repository.SessionRepository {
	return &sessionRepository{}
}

// ErrorMsg generates and returns error message.
func (repo *sessionRepository) ErrorMsg(method model.RepositoryMethod, err error) error {
	ex := &model.RepositoryError{
		BaseErr:          err,
		RepositoryMethod: method,
		DomainModelName:  model.DomainModelNameSession,
	}

	return ex
}

// GetSessionByID gets and returns a record specified by id.
func (repo *sessionRepository) GetSessionByID(ctx context.Context, m query.SQLManager, id string) (*model.Session, error) {
	q := "SELECT id, user_id, created_at FROM sessions WHERE id=?"

	list, err := repo.list(ctx, m, model.RepositoryMethodREAD, q, id)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:         err,
			PropertyName:    model.IDProperty,
			PropertyValue:   id,
			DomainModelName: model.DomainModelNameSession,
		}
		return nil, errors.Wrapf(err, "session data is 0")
	}

	if err != nil {
		err = errors.Wrap(err, "failed to list session")
		return nil, repo.ErrorMsg(model.RepositoryMethodREAD, err)
	}

	return list[0], nil
}

// list gets and returns list of records.
func (repo *sessionRepository) list(ctx context.Context, m query.SQLManager, method model.RepositoryMethod, q string, args ...interface{}) (sessions []*model.Session, err error) {
	stmt, err := m.PrepareContext(ctx, q)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare context")
		return nil, repo.ErrorMsg(method, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		err = errors.Wrap(err, "failed to execute context")
		return nil, repo.ErrorMsg(method, err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			logger.Logger.Error("rows.Close", zap.String("error message", err.Error()))
		}
	}()

	list := make([]*model.Session, 0)
	for rows.Next() {
		session := &model.Session{}

		err = rows.Scan(
			&session.ID,
			&session.UserID,
			&session.CreatedAt,
		)

		if err != nil {
			err = errors.Wrap(err, "failed to scan rows")
			return nil, repo.ErrorMsg(method, err)
		}

		list = append(list, session)
	}

	return list, nil
}

// InsertSession insert a record.
func (repo *sessionRepository) InsertSession(ctx context.Context, m query.SQLManager, session *model.Session) error {
	q := "INSERT INTO sessions (id, user_id, created_at) VALUES (?, ?, NOW())"
	stmt, err := m.PrepareContext(ctx, q)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare context")
		return repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	result, err := stmt.ExecContext(ctx, session.ID, session.UserID)
	if err != nil {
		err = errors.Wrap(err, "failed to execute context")
		return repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = errors.Errorf("total affected: %d ", affect)
		return repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	return nil
}

// DeleteSession delete a record.
func (repo *sessionRepository) DeleteSession(ctx context.Context, m query.SQLManager, id string) error {
	q := "DELETE FROM sessions WHERE id=?"

	stmt, err := m.PrepareContext(ctx, q)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare context")
		return repo.ErrorMsg(model.RepositoryMethodDELETE, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		err = errors.Wrap(err, "failed to execute context")
		return repo.ErrorMsg(model.RepositoryMethodDELETE, err)
	}

	affect, err := result.RowsAffected()
	if err != nil {
		err = errors.Wrap(err, "failed to get rows affected")
		return repo.ErrorMsg(model.RepositoryMethodDELETE, err)
	}
	if affect != 1 {
		err = errors.Errorf("total affected: %d ", affect)
		return repo.ErrorMsg(model.RepositoryMethodDELETE, err)
	}

	return nil
}
