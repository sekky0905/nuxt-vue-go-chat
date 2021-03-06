package db

import (
	"context"

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
		BaseErr:          err,
		RepositoryMethod: method,
		DomainModelName:  model.DomainModelNameThread,
	}
}

// ListThreads lists ThreadList.
func (repo *threadRepository) ListThreads(ctx context.Context, m query.SQLManager, cursor uint32, limit int) (*model.ThreadList, error) {
	q := `SELECT t.id, t.title, u.id, u.name, t.created_at, t.updated_at
	FROM threads AS t
	INNER JOIN users AS u
	ON t.user_id = u.id
	WHERE t.id >= ?
	ORDER BY t.id ASC
	LIMIT ?;`

	limitForCheckHasNext := readyLimitForHasNext(limit)
	threads, err := repo.list(ctx, m, model.RepositoryMethodREAD, q, cursor, limitForCheckHasNext)

	length := len(threads)

	if length == 0 {
		err = &model.NoSuchDataError{
			BaseErr:         err,
			DomainModelName: model.DomainModelNameThread,
		}
		return nil, err
	}

	if err != nil {
		err = errors.Wrap(err, "failed to list threads")
		return nil, repo.ErrorMsg(model.RepositoryMethodLIST, err)
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
	q := `SELECT t.id, t.title, u.id, u.name, t.created_at, t.updated_at
	FROM threads AS t
	INNER JOIN users AS u
	ON t.user_id = u.id
	WHERE t.id=?
	LIMIT 1;`

	list, err := repo.list(ctx, m, model.RepositoryMethodREAD, q, id)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:         err,
			PropertyName:    model.IDProperty,
			PropertyValue:   id,
			DomainModelName: model.DomainModelNameThread,
		}
		return nil, err
	}

	if err != nil {
		err = errors.Wrap(err, "failed to list threads")
		return nil, repo.ErrorMsg(model.RepositoryMethodLIST, err)
	}

	return list[0], nil
}

// GetThreadByTitle gets and returns a record specified by title.
func (repo *threadRepository) GetThreadByTitle(ctx context.Context, m query.SQLManager, name string) (*model.Thread, error) {
	q := `SELECT t.id, t.title, u.id, u.name, t.created_at, t.updated_at
	FROM threads AS t
	INNER JOIN users AS u
	ON t.user_id = u.id
	WHERE t.title=?
	LIMIT 1;`

	list, err := repo.list(ctx, m, model.RepositoryMethodREAD, q, name)

	if len(list) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:         err,
			PropertyName:    model.NameProperty,
			PropertyValue:   name,
			DomainModelName: model.DomainModelNameThread,
		}
		return nil, err
	}

	if err != nil {
		err = errors.Wrap(err, "failed to list threads")
		return nil, repo.ErrorMsg(model.RepositoryMethodLIST, err)
	}

	return list[0], nil
}

// list gets and returns list of records.
func (repo *threadRepository) list(ctx context.Context, m query.SQLManager, method model.RepositoryMethod, q string, args ...interface{}) (threads []*model.Thread, err error) {
	stmt, err := m.PrepareContext(ctx, q)
	if err != nil {
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
		err = errors.Wrap(err, "failed to query context")
		return nil, repo.ErrorMsg(method, err)
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
			err = errors.Wrap(err, "failed to scan rows")
			return nil, repo.ErrorMsg(method, err)
		}

		list = append(list, thread)
	}

	return list, nil
}

// InsertThread insert a record.
func (repo *threadRepository) InsertThread(ctx context.Context, m query.SQLManager, thread *model.Thread) (uint32, error) {
	q := "INSERT INTO threads (title, user_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())"
	stmt, err := m.PrepareContext(ctx, q)
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

	result, err := stmt.ExecContext(ctx, thread.Title, thread.User.ID)
	if err != nil {
		err = errors.Wrap(err, "failed to execute context")
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = errors.Errorf("total affected is: %d", affect)
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		err = errors.Wrap(err, "failed to get last insert id")
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	return uint32(id), nil
}

// UpdateThread updates a record.
func (repo *threadRepository) UpdateThread(ctx context.Context, m query.SQLManager, id uint32, thread *model.Thread) error {
	q := "UPDATE threads SET title=?, updated_at=NOW() WHERE id=?"
	stmt, err := m.PrepareContext(ctx, q)
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

	result, err := stmt.ExecContext(ctx, thread.Title, id)
	if err != nil {
		err = errors.Wrap(err, "failed to execute context")
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, err)
	}

	affect, err := result.RowsAffected()
	if err != nil {
		err = errors.Wrap(err, "failed to get rows affected")
		return repo.ErrorMsg(model.RepositoryMethodDELETE, err)
	}

	if affect != 1 {
		err = errors.Errorf("total affected is: %d", affect)
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, err)
	}

	return nil
}

// DeleteThread delete a record.
func (repo *threadRepository) DeleteThread(ctx context.Context, m query.SQLManager, id uint32) error {
	q := "DELETE FROM threads WHERE id=?"

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
		err = errors.Errorf("total affected is: %d", affect)
		return repo.ErrorMsg(model.RepositoryMethodDELETE, err)
	}

	return nil
}
