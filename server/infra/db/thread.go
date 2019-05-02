package db

import (
	"context"
	"fmt"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"
	"go.uber.org/zap"
)

// threadRepository is repository of thread.
type threadRepository struct {
}

// NewThreadRepository generates and returns ThreadRepository.
func NewThreadRepository() repository.ThreadRepository {
	return &threadRepository{}
}

// ErrorMsg generates and returns error message.
func (repo *threadRepository) ErrorMsg(method model.RepositoryMethod, err error) error {
	return &model.RepositoryError{
		BaseErr:                     err,
		RepositoryMethod:            method,
		DomainModelNameForDeveloper: model.DomainModelNameThreadForDeveloper,
		DomainModelNameForUser:      model.DomainModelNameThreadForUser,
	}
}

// ListThreads lists ThreadList.
func (repo *threadRepository) ListThreads(ctx context.Context, m query.SQLManager, cursor uint32, limit int) (*model.ThreadList, error) {
	query := `SELECT t.id, t.title, u.id, u.name, t.created_at, t.updated_at
	FROM threads AS t
	INNER JOIN users AS u
	ON t.user_id = u.id
	WHERE t.id >= ?
	ORDER BY t.id ASC
	LIMIT ?;`

	limitForCheckHasNext := readyLimitForHasNext(limit)
	threads, err := repo.list(ctx, m, model.RepositoryMethodREAD, query, cursor, limitForCheckHasNext)

	length := len(threads)

	if length == 0 {
		err = &model.NoSuchDataError{
			BaseErr:                     err,
			DomainModelNameForDeveloper: model.DomainModelNameThreadForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameThreadForUser,
		}
		return nil, err
	}

	if err != nil {
		return nil, repo.ErrorMsg(model.RepositoryMethodLIST, errors.WithStack(err))
	}

	hasNext := checkHasNext(length, limit)
	if hasNext {
		cursor = threads[limitForCheckHasNext-1].ID
	} else {
		cursor = 0
	}

	if length == limitForCheckHasNext {
		// exclude thread for cursor
		return &model.ThreadList{Threads: threads[:limitForCheckHasNext-1], HasNext: hasNext, Cursor: cursor}, nil
	}

	return &model.ThreadList{Threads: threads, HasNext: hasNext, Cursor: cursor}, nil
}

// GetThreadByID gets and returns a record specified by id.
func (repo *threadRepository) GetThreadByID(ctx context.Context, m query.SQLManager, id uint32) (*model.Thread, error) {
	query := `SELECT t.id, t.title, u.id, u.name, t.created_at, t.updated_at
	FROM threads AS t
	INNER JOIN users AS u
	ON t.user_id = u.id
	WHERE t.id=?
	LIMIT 1;`

	list, err := repo.list(ctx, m, model.RepositoryMethodREAD, query, id)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:                     err,
			PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
			PropertyNameForUser:         model.IDPropertyForUser,
			PropertyValue:               id,
			DomainModelNameForDeveloper: model.DomainModelNameThreadForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameThreadForUser,
		}
		return nil, err
	}

	if err != nil {
		return nil, repo.ErrorMsg(model.RepositoryMethodLIST, errors.WithStack(err))
	}

	return list[0], nil
}

// GetThreadByTitle gets and returns a record specified by title.
func (repo *threadRepository) GetThreadByTitle(ctx context.Context, m query.SQLManager, name string) (*model.Thread, error) {
	query := `SELECT t.id, t.title, u.id, u.name, t.created_at, t.updated_at
	FROM threads AS t
	INNER JOIN users AS u
	ON t.user_id = u.id
	WHERE t.title=?
	LIMIT 1;`

	list, err := repo.list(ctx, m, model.RepositoryMethodREAD, query, name)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:                     err,
			PropertyNameForDeveloper:    model.NamePropertyForDeveloper,
			PropertyNameForUser:         model.NamePropertyForUser,
			PropertyValue:               name,
			DomainModelNameForDeveloper: model.DomainModelNameThreadForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameThreadForUser,
		}
		return nil, err
	}

	if err != nil {
		return nil, repo.ErrorMsg(model.RepositoryMethodLIST, errors.WithStack(err))
	}

	return list[0], nil
}

// list gets and returns list of records.
func (repo *threadRepository) list(ctx context.Context, m query.SQLManager, method model.RepositoryMethod, query string, args ...interface{}) (threads []*model.Thread, err error) {
	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.WithStack(repo.ErrorMsg(method, err))
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, repo.ErrorMsg(method, errors.WithStack(err))
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			logger.Logger.Error("rows.Close", zap.String("error message", err.Error()))
		}
	}()

	list := make([]*model.Thread, 0)
	for rows.Next() {
		thread := &model.Thread{
			User: &model.User{},
		}

		err = rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.User.ID,
			&thread.User.Name,
			&thread.CreatedAt,
			&thread.UpdatedAt,
		)

		if err != nil {
			return nil, repo.ErrorMsg(method, errors.WithStack(err))
		}

		list = append(list, thread)
	}

	return list, nil
}

// InsertThread insert a record.
func (repo *threadRepository) InsertThread(ctx context.Context, m query.SQLManager, thread *model.Thread) (uint32, error) {
	query := "INSERT INTO threads (title, user_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())"
	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		return model.InvalidID, errors.WithStack(repo.ErrorMsg(model.RepositoryMethodInsert, err))
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	result, err := stmt.ExecContext(ctx, thread.Title, thread.User.ID)
	if err != nil {
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, errors.WithStack(err))
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = fmt.Errorf("total affected is: %d", affect)
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, errors.WithStack(err))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, errors.WithStack(err))
	}

	return uint32(id), nil
}

// UpdateThread updates a record.
func (repo *threadRepository) UpdateThread(ctx context.Context, m query.SQLManager, id uint32, thread *model.Thread) error {
	query := "UPDATE threads SET title=?, updated_at=NOW() WHERE id=?"
	stmt, err := m.PrepareContext(ctx, query)
	if err != nil {
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, errors.WithStack(err))
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	result, err := stmt.ExecContext(ctx, thread.Title, id)
	if err != nil {
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, errors.WithStack(err))
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = fmt.Errorf("total affected is: %d", affect)
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, errors.WithStack(err))
	}

	return nil
}

// DeleteThread delete a record.
func (repo *threadRepository) DeleteThread(ctx context.Context, m query.SQLManager, id uint32) error {
	query := "DELETE FROM threads WHERE id=?"

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
		err = fmt.Errorf("total affected is: %d", affect)
		return repo.ErrorMsg(model.RepositoryMethodDELETE, errors.WithStack(err))
	}

	return nil
}
