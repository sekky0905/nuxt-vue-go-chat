package db

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"
)

// userRepository is repository of user.
type userRepository struct {
}

// NewUserRepository generates and returns userRepository.
func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

// ErrorMsg generates and returns error message.
func (repo *userRepository) ErrorMsg(method model.RepositoryMethod, err error) error {
	return &model.RepositoryError{
		BaseErr:          err,
		RepositoryMethod: method,
		DomainModelName:  model.DomainModelNameUser,
	}
}

// GetUserByID gets and returns a record specified by id.
func (repo *userRepository) GetUserByID(ctx context.Context, m query.SQLManager, id uint32) (*model.User, error) {
	query := "SELECT id, name, session_id, password, created_at, updated_at FROM users WHERE id=?"

	list, err := repo.list(ctx, m, model.RepositoryMethodREAD, query, id)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:         err,
			PropertyName:    model.IDProperty,
			PropertyValue:   id,
			DomainModelName: model.DomainModelNameUser,
		}
		return nil, err
	}

	if err != nil {
		err = errors.Wrap(err, "failed to list users")
		return nil, repo.ErrorMsg(model.RepositoryMethodREAD, err)
	}

	return list[0], nil
}

// GetUserByName gets and returns a record specified by name.
func (repo *userRepository) GetUserByName(ctx context.Context, m query.SQLManager, name string) (*model.User, error) {
	query := "SELECT id, name, session_id, password, created_at, updated_at FROM users WHERE name=?"
	list, err := repo.list(ctx, m, model.RepositoryMethodREAD, query, name)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:         err,
			PropertyName:    model.NameProperty,
			PropertyValue:   name,
			DomainModelName: model.DomainModelNameUser,
		}
		return nil, err
	}

	if err != nil {
		err = errors.Wrap(err, "failed to list users")
		return nil, repo.ErrorMsg(model.RepositoryMethodREAD, err)
	}

	return list[0], nil
}

// list gets and returns list of records.
func (repo *userRepository) list(ctx context.Context, m query.SQLManager, method model.RepositoryMethod, query string, args ...interface{}) (users []*model.User, err error) {
	stmt, err := m.PrepareContext(ctx, query)
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

	list := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}

		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.SessionID,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			err = errors.Wrap(err, "failed to scan rows")
			return nil, repo.ErrorMsg(method, err)
		}

		list = append(list, user)
	}

	return list, nil
}

// InsertUser insert a record.
func (repo *userRepository) InsertUser(ctx context.Context, m query.SQLManager, user *model.User) (uint32, error) {
	query := "INSERT INTO users (name, session_id, password, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"
	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare context")
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	result, err := stmt.ExecContext(ctx, user.Name, user.SessionID, user.Password)
	if err != nil {
		err = errors.Wrap(err, "failed to execute context")
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	affect, err := result.RowsAffected()
	if err != nil {
		err = errors.Wrap(err, "failed to get rows affected")
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	if affect != 1 {
		err = errors.Errorf("total affected: %d ", affect)
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	return uint32(id), nil
}

// UpdateUser updates a record.
func (repo *userRepository) UpdateUser(ctx context.Context, m query.SQLManager, id uint32, user *model.User) error {
	query := "UPDATE users SET session_id=?, password=?, updated_at=NOW() WHERE id=?"

	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare context")
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, err)
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	result, err := stmt.ExecContext(ctx, user.SessionID, user.Password, id)
	if err != nil {
		err = errors.Wrap(err, "failed to execute context")
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, err)
	}

	affect, err := result.RowsAffected()
	if err != nil {
		err = errors.Wrap(err, "failed to get rows affected")
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, err)
	}
	if affect != 1 {
		err = errors.Errorf("total affected: %d ", affect)
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, err)
	}

	return nil
}

// DeleteUser delete a record.
func (repo *userRepository) DeleteUser(ctx context.Context, m query.SQLManager, id uint32) error {
	query := "DELETE FROM users WHERE id=?"

	stmt, err := m.PrepareContext(ctx, query)
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
