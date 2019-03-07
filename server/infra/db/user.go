package db

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	. "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// userRepository is repository of user.
type userRepository struct {
	ctx context.Context
}

// NewUserRepository は、userRepository生成する。
func NewUserRepository(ctx context.Context) UserRepository {
	return &userRepository{
		ctx: ctx,
	}
}

// ErrorMsg generates and returns error message.
func (repo *userRepository) ErrorMsg(method model.RepositoryMethod, err error) error {
	return &model.RepositoryError{
		BaseErr:                     err,
		RepositoryMethod:            method,
		DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
		DomainModelNameForUser:      model.DomainModelNameUserForUser,
	}
}

// GetUserByID gets and returns a record specified by id.
func (repo *userRepository) GetUserByID(m DBManager, id uint32) (*model.User, error) {
	query := "SELECT id, name, session_id, password, created_at, updated_at FROM users WHERE id=?"

	list, err := repo.list(m, model.RepositoryMethodREAD, query, id)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:                     err,
			PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
			PropertyNameForUser:         model.IDPropertyForUser,
			PropertyValue:               id,
			DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameUserForUser,
		}
		return nil, err
	}

	if err != nil {
		return nil, errors.WithStack(repo.ErrorMsg(model.RepositoryMethodREAD, errors.WithStack(err)))
	}

	return list[0], nil
}

// GetUserByName gets and returns a record specified by name.
func (repo *userRepository) GetUserByName(m DBManager, name string) (*model.User, error) {
	query := "SELECT id, name, session_id, password, created_at, updated_at FROM users WHERE name=?"
	list, err := repo.list(m, model.RepositoryMethodREAD, query, name)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:                     err,
			PropertyNameForDeveloper:    model.NamePropertyForDeveloper,
			PropertyNameForUser:         model.NamePropertyForUser,
			PropertyValue:               name,
			DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameUserForUser,
		}
		return nil, err
	}

	if err != nil {
		return nil, errors.WithStack(repo.ErrorMsg(model.RepositoryMethodREAD, err))
	}

	return list[0], nil
}

// list gets and returns list of records.
func (repo *userRepository) list(m DBManager, method model.RepositoryMethod, query string, args ...interface{}) (users []*model.User, err error) {
	stmt, err := m.PrepareContext(repo.ctx, query)
	if err != nil {
		return nil, repo.ErrorMsg(method, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	rows, err := stmt.QueryContext(repo.ctx, args...)
	if err != nil {
		return nil, errors.WithStack(repo.ErrorMsg(method, err))
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Error(err.Error())
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
			return nil, errors.WithStack(repo.ErrorMsg(method, err))
		}

		list = append(list, user)
	}

	return list, nil
}

// InsertUser insert a record.
func (repo *userRepository) InsertUser(m DBManager, user *model.User) (uint32, error) {
	query := "INSERT INTO users (name, session_id, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	stmt, err := m.PrepareContext(repo.ctx, query)
	if err != nil {
		return model.InvalidID, errors.WithStack(repo.ErrorMsg(model.RepositoryMethodInsert, err))
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	result, err := stmt.ExecContext(repo.ctx, user.Name, user.SessionID, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return model.InvalidID, errors.WithStack(repo.ErrorMsg(model.RepositoryMethodInsert, err))
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = fmt.Errorf("total affected: %d ", affect)
		return model.InvalidID, errors.WithStack(repo.ErrorMsg(model.RepositoryMethodInsert, err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.InvalidID, errors.WithStack(repo.ErrorMsg(model.RepositoryMethodInsert, err))
	}

	return uint32(id), nil
}

// UpdateUser updates a record.
func (repo *userRepository) UpdateUser(m DBManager, id uint32, user *model.User) error {
	query := "UPDATE users SET session_id=?, password=?, updated_at=? WHERE id=?"

	stmt, err := m.PrepareContext(repo.ctx, query)
	if err != nil {
		return errors.WithStack(repo.ErrorMsg(model.RepositoryMethodUPDATE, err))
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	result, err := stmt.ExecContext(repo.ctx, user.SessionID, user.Password, user.UpdatedAt, id)
	if err != nil {
		return errors.WithStack(repo.ErrorMsg(model.RepositoryMethodUPDATE, err))
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = fmt.Errorf("total affected: %d ", affect)
		return errors.WithStack(repo.ErrorMsg(model.RepositoryMethodUPDATE, err))
	}

	return nil
}
func (repo *userRepository) DeleteUser(m DBManager, id uint32) error { return nil }
