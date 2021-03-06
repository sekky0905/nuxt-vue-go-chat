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

// commentRepository is repository of comment.
type commentRepository struct {
}

// NewCommentRepository generates and returns CommentRepository.
func NewCommentRepository() repository.CommentRepository {
	return &commentRepository{}
}

// ErrorMsg generates and returns error message.
func (repo *commentRepository) ErrorMsg(method model.RepositoryMethod, err error) error {
	return &model.RepositoryError{
		BaseErr:          err,
		RepositoryMethod: method,
		DomainModelName:  model.DomainModelNameComment,
	}
}

// ListThreads lists ThreadList.
func (repo *commentRepository) ListComments(ctx context.Context, m query.SQLManager, threadID uint32, limit int, cursor uint32) (*model.CommentList, error) {
	q := `SELECT c.id, c.content, u.id, u.name, c.thread_id, c.created_at, c.updated_at
	FROM comments AS c
	INNER JOIN users AS u
	ON c.user_id = u.id
	WHERE c.id >= ?
	AND c.thread_id = ?
	ORDER BY c.id ASC
	LIMIT ?;`

	limitForCheckHasNext := readyLimitForHasNext(limit)
	comments, err := repo.list(ctx, m, model.RepositoryMethodREAD, q, cursor, threadID, limitForCheckHasNext)

	length := len(comments)

	if length == 0 {
		err = &model.NoSuchDataError{
			BaseErr:         err,
			DomainModelName: model.DomainModelNameComment,
		}
		return nil, err
	}

	if err != nil {
		err = errors.Wrap(err, "failed to list comments")
		return nil, repo.ErrorMsg(model.RepositoryMethodLIST, err)
	}

	hasNext := checkHasNext(length, limit)
	if hasNext {
		cursor = comments[limitForCheckHasNext-1].ID
	} else {
		cursor = 0
	}

	if length == limitForCheckHasNext {
		// exclude thread for cursor
		return &model.CommentList{
			Comments: comments[:limitForCheckHasNext-1],
			HasNext:  hasNext,
			Cursor:   cursor,
		}, nil
	}

	return &model.CommentList{
		Comments: comments,
		HasNext:  hasNext,
		Cursor:   cursor,
	}, nil
}

// GetThreadByID gets and returns a record specified by id.
func (repo *commentRepository) GetCommentByID(ctx context.Context, m query.SQLManager, id uint32) (*model.Comment, error) {
	q := `SELECT c.id, c.content, u.id, u.name, c.thread_id, c.created_at, c.updated_at
	FROM comments AS c
	INNER JOIN users AS u
	ON c.user_id = u.id
	WHERE c.id=?
	LIMIT 1;`

	comments, err := repo.list(ctx, m, model.RepositoryMethodREAD, q, id)

	if len(comments) == 0 {
		err = &model.NoSuchDataError{
			BaseErr:         err,
			PropertyName:    model.IDProperty,
			PropertyValue:   id,
			DomainModelName: model.DomainModelNameComment,
		}
		return nil, err
	}

	if err != nil {
		err = errors.Wrap(err, "failed to list comments")
		return nil, repo.ErrorMsg(model.RepositoryMethodREAD, err)
	}

	return comments[0], nil
}

// list gets and returns list of records.
func (repo *commentRepository) list(ctx context.Context, m query.SQLManager, method model.RepositoryMethod, q string, args ...interface{}) (comments []*model.Comment, err error) {
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

	list := make([]*model.Comment, 0)
	for rows.Next() {
		comment := &model.Comment{
			User: &model.User{},
		}

		err = rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.User.ID,
			&comment.User.Name,
			&comment.ThreadID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)

		if err != nil {
			err = errors.Wrap(err, "failed to scan rows")
			return nil, repo.ErrorMsg(method, err)
		}

		list = append(list, comment)
	}

	return list, nil
}

// InsertThread insert a record.
func (repo *commentRepository) InsertComment(ctx context.Context, m query.SQLManager, comment *model.Comment) (uint32, error) {
	q := "INSERT INTO comments (content, user_id, thread_id, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW());"
	stmt, err := m.PrepareContext(ctx, q)
	if err != nil {
		err = errors.Wrap(err, "failed to list prepare context")
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error("stmt.Close", zap.String("error message", err.Error()))
		}
	}()

	result, err := stmt.ExecContext(ctx, comment.Content, comment.User.ID, comment.ThreadID)
	if err != nil {
		err = errors.Wrap(err, "failed to list execute context")
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = errors.Errorf("total affected id: %d ", affect)
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		err = errors.Wrap(err, "failed to get last insert id")
		return model.InvalidID, repo.ErrorMsg(model.RepositoryMethodInsert, err)
	}

	return uint32(id), nil
}

// UpdateComment updates a record.
func (repo *commentRepository) UpdateComment(ctx context.Context, m query.SQLManager, id uint32, comment *model.Comment) error {
	q := "UPDATE comments SET content=?, updated_at= NOW() WHERE id=?;"

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

	result, err := stmt.ExecContext(ctx, comment.Content, id)
	if err != nil {
		err = errors.Wrap(err, "failed to execute context")
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, err)
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = errors.Errorf("total affected id: %d ", affect)
		return repo.ErrorMsg(model.RepositoryMethodUPDATE, err)
	}

	return nil
}

// DeleteComment delete a record.
func (repo *commentRepository) DeleteComment(ctx context.Context, m query.SQLManager, id uint32) error {
	q := "DELETE FROM comments WHERE id=?;"

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
		err = errors.Wrap(err, "failed get rows affected")
		return repo.ErrorMsg(model.RepositoryMethodDELETE, err)
	}
	if affect != 1 {
		err = errors.Errorf("total affected id: %d ", affect)
		return repo.ErrorMsg(model.RepositoryMethodDELETE, err)
	}

	return nil
}
