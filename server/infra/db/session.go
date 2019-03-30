package db

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
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
		BaseErr:                     err,
		RepositoryMethod:            method,
		DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
		DomainModelNameForUser:      model.DomainModelNameSessionForUser,
	}

	return ex
}

// GetSessionByID gets and returns a record specified by id.
func (repo *sessionRepository) GetSessionByID(ctx context.Context, m repository.SQLManager, id string) (*model.Session, error) {
	query := "SELECT id, user_id, created_at FROM sessions WHERE id=?"

	list, err := repo.list(ctx, m, model.RepositoryMethodREAD, query, id)

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
func (repo *sessionRepository) list(ctx context.Context, m repository.SQLManager, method model.RepositoryMethod, query string, args ...interface{}) (sessions []*model.Session, err error) {
	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		return nil, repo.ErrorMsg(method, errors.WithStack(err))
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		err = repo.ErrorMsg(method, errors.WithStack(err))
		return nil, err
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
			return nil, repo.ErrorMsg(method, errors.WithStack(err))
		}

		list = append(list, session)
	}

	return list, nil
}

// InsertSession insert a record.
func (repo *sessionRepository) InsertSession(ctx context.Context, m repository.SQLManager, session *model.Session) error {
	query := "INSERT INTO sessions (id, user_id, created_at) VALUES (?, ?, ?)"
	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		return errors.WithStack(repo.ErrorMsg(model.RepositoryMethodInsert, err))
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	result, err := stmt.ExecContext(ctx, session.ID, session.UserID, session.CreatedAt)
	if err != nil {
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
func (repo *sessionRepository) DeleteSession(ctx context.Context, m repository.SQLManager, id string) error {
	query := "DELETE FROM sessions WHERE id=?"

	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		return repo.ErrorMsg(model.RepositoryMethodDELETE, errors.WithStack(err))
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return repo.ErrorMsg(model.RepositoryMethodDELETE, errors.WithStack(err))
	}

	affect, err := result.RowsAffected()
	if err != nil {
		return repo.ErrorMsg(model.RepositoryMethodDELETE, errors.WithStack(err))
	}
	if affect != 1 {
		err = fmt.Errorf("total affected: %d ", affect)
		return repo.ErrorMsg(model.RepositoryMethodDELETE, errors.WithStack(err))
	}

	return nil
}
