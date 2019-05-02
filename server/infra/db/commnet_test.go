package db

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewCommentRepository(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want repository.CommentRepository
	}{
		{
			name: "When given appropriate args, returns commentRepository",
			args: args{
				ctx: context.Background(),
			},
			want: &commentRepository{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommentRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommentRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commentRepository_ErrorMsg(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		method model.RepositoryMethod
		err    error
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr *model.RepositoryError
	}{
		{
			name: "When given appropriate args, returns appropriate error",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				method: model.RepositoryMethodInsert,
				err:    errors.New(model.ErrorMessageForTest),
			},
			wantErr: &model.RepositoryError{
				BaseErr:                     errors.New(model.ErrorMessageForTest),
				RepositoryMethod:            model.RepositoryMethodInsert,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &commentRepository{}
			if err := repo.ErrorMsg(tt.args.method, tt.args.err); errors.Cause(err).Error() != tt.wantErr.Error() {
				t.Errorf("commentRepository{ErrorMsg() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}

func Test_commentRepository_ListComments(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		return
	}

	defer db.Close()

	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx       context.Context
		m         query.SQLManager
		commentID uint32
		limit     int
		cursor    uint32
	}

	tests := []struct {
		name       string
		repo       *commentRepository
		args       args
		want       *model.CommentList
		returnMock []*model.Comment
		wantErr    error
	}{
		{
			name: "When limit = 20, cursor = 1 are given and there are over 21 data, ListComments returns CommentList which has Comments(ID: 1~20), HasNext = yes, Cursor = 21",
			repo: &commentRepository{},
			args: args{
				ctx:       context.Background(),
				m:         db,
				commentID: model.ThreadValidIDForTest,
				limit:     20,
				cursor:    1,
			},
			want: &model.CommentList{
				Comments: testutil.GenerateCommentHelper(1, 20),
				HasNext:  true,
				Cursor:   21,
			},
			returnMock: testutil.GenerateCommentHelper(1, 21),
			wantErr:    nil,
		},
		{
			name: "When limit = 20, cursor = 21 are given and there are over 41 data, ListComments returns CommentList which has 21 Comments(ID: 21~40), HasNext = yes, Cursor = 41",
			repo: &commentRepository{},
			args: args{
				ctx:       context.Background(),
				m:         db,
				commentID: model.ThreadValidIDForTest,
				limit:     20,
				cursor:    21,
			},
			want: &model.CommentList{
				Comments: testutil.GenerateCommentHelper(21, 40),
				HasNext:  true,
				Cursor:   41,
			},
			returnMock: testutil.GenerateCommentHelper(21, 41),
			wantErr:    nil,
		},
		{
			name: "When limit = 20, cursor = 1 are given and there are over 10 data, ListComments returns CommentList which has 10 Comments(ID: 1~10), HasNext = yes, Cursor = 10",
			repo: &commentRepository{},
			args: args{
				ctx:       context.Background(),
				m:         db,
				commentID: model.ThreadValidIDForTest,
				limit:     20,
				cursor:    1,
			},
			want: &model.CommentList{
				Comments: testutil.GenerateCommentHelper(1, 10),
				HasNext:  false,
				Cursor:   0,
			},
			returnMock: testutil.GenerateCommentHelper(1, 10),
			wantErr:    nil,
		},
		{
			name: "When limit = 20, cursor = 1 are given and there are no data, ListComments returns error",
			repo: &commentRepository{},
			args: args{
				ctx:       context.Background(),
				m:         db,
				commentID: model.ThreadValidIDForTest,
				limit:     20,
				cursor:    1,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				BaseErr:                     err,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `SELECT (.+)
	FROM comments AS c
	INNER JOIN users AS u
	(.+);`
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WithArgs(tt.args.cursor, readyLimitForHasNext(tt.args.limit)).WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"c.id", "c.content", "u.id", "u.name", "c.comment_id", "c.created_at", "c.updated_at"})

				for _, comment := range tt.returnMock {
					rows.AddRow(comment.ID, comment.Content, comment.User.ID, comment.User.Name, comment.ThreadID, comment.CreatedAt, comment.UpdatedAt)
				}

				prep.ExpectQuery().WithArgs(tt.args.cursor, tt.args.commentID, readyLimitForHasNext(tt.args.limit)).WillReturnRows(rows)
			}

			repo := &commentRepository{}
			got, err := repo.ListComments(tt.args.ctx, tt.args.m, tt.args.commentID, tt.args.limit, tt.args.cursor)
			if tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("commentRepository.ListComments() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commentepository.ListComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commentRepository_GetCommentByID(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		return
	}

	defer db.Close()

	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx context.Context
		m   query.SQLManager
		id  uint32
	}

	tests := []struct {
		name    string
		args    args
		want    *model.Comment
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When a comment specified by id exists, returns a comment",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentValidIDForTest,
			},
			want: &model.Comment{
				ID:       model.CommentValidIDForTest,
				ThreadID: model.ThreadValidIDForTest,
				User: &model.User{
					ID:   model.UserValidIDForTest,
					Name: model.UserNameForTest,
				},
				Content: model.CommentContentForTest,
			},
			wantErr: nil,
		},
		{
			name: "When a comment specified by id does not exist, returns NoSuchDataError",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentInValidIDForTest,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
				PropertyNameForUser:         model.IDPropertyForUser,
				PropertyValue:               model.UserInValidIDForTest,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `SELECT (.+)
	FROM comments AS c
	INNER JOIN users AS u
	(.+);`
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"c.id", "c.content", "u.id", "u.name", "c.comment_id", "c.created_at", "c.updated_at"}).
					AddRow(tt.want.ID, tt.want.Content, tt.want.User.ID, tt.want.User.Name, tt.want.ThreadID, tt.want.CreatedAt, tt.want.UpdatedAt)
				prep.ExpectQuery().WithArgs(tt.want.ID).WillReturnRows(rows)
			}

			repo := &commentRepository{}
			got, err := repo.GetCommentByID(tt.args.ctx, tt.args.m, tt.args.id)

			if tt.wantErr != nil {
				if !reflect.DeepEqual(err, tt.wantErr) {
					t.Errorf("commentRepository.GetCommentByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commentRepository.GetCommentByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commentRepository_InsertComment(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx     context.Context
		m       query.SQLManager
		comment *model.Comment
		err     error
	}

	tests := []struct {
		name        string
		args        args
		rowAffected int64
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a comment which has ID, ThreadID, User, Content is given, returns ID",
			args: args{
				ctx: context.Background(),
				m:   db,
				comment: &model.Comment{
					ID:       model.CommentValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
			},
			rowAffected: 1,
			wantErr:     nil,
		},
		{
			name: "when RowAffected is 0、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				comment: &model.Comment{
					ID:       model.CommentInValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodInsert,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				comment: &model.Comment{
					ID:       model.CommentInValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodInsert,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				comment: &model.Comment{
					ID:       model.CommentInValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
				err: errors.New(model.ErrorMessageForTest),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodInsert,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "INSERT INTO comments"
			prep := mock.ExpectPrepare(query)

			exec := prep.ExpectExec().WithArgs(tt.args.comment.Content, tt.args.comment.User.ID, tt.args.comment.ThreadID, tt.args.comment.CreatedAt, tt.args.comment.UpdatedAt)

			if tt.args.err != nil {
				exec.WillReturnError(tt.args.err)
			} else {
				exec.WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &commentRepository{}

			_, err := repo.InsertComment(tt.args.ctx, tt.args.m, tt.args.comment)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("commentRepository.InsertComment() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func Test_commentRepository_UpdateComment(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx     context.Context
		m       query.SQLManager
		id      uint32
		comment *model.Comment
		err     error
	}

	tests := []struct {
		name        string
		args        args
		rowAffected int64
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a comment which has ID, ThreadID, User, Content is given, returns nil",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentValidIDForTest,
				comment: &model.Comment{
					ID:       model.CommentInValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
			},
			rowAffected: 1,
			wantErr:     nil,
		},
		{
			name: "when RowAffected is 0、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentInValidIDForTest,
				comment: &model.Comment{
					ID:       model.CommentInValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodUPDATE,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentInValidIDForTest,
				comment: &model.Comment{
					ID:       model.CommentInValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodUPDATE,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentInValidIDForTest,
				comment: &model.Comment{
					ID:       model.CommentInValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
				err: errors.New(model.ErrorMessageForTest),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodUPDATE,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "UPDATE comments SET title=\\?, updated_at=\\? WHERE id=\\?"
			prep := mock.ExpectPrepare(query)

			exec := prep.ExpectExec().WithArgs(tt.args.comment.Content, tt.args.comment.User.ID, tt.args.comment.ThreadID, tt.args.comment.CreatedAt, tt.args.comment.UpdatedAt)

			if tt.args.err != nil {
				exec.WillReturnError(tt.args.err)
			} else {
				exec.WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &commentRepository{}
			err := repo.UpdateComment(tt.args.ctx, tt.args.m, tt.args.id, tt.args.comment)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("commentRepository.UpdateComment() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func Test_commentRepository_DeleteComment(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx context.Context
		m   query.SQLManager
		id  uint32
		err error
	}

	tests := []struct {
		name        string
		rowAffected int64
		args        args
		wantErr     *model.RepositoryError
	}{
		{
			name:        "When a comment specified by id exists, returns nil",
			rowAffected: 1,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentValidIDForTest,
			},
			wantErr: nil,
		},
		{
			name:        "when RowAffected is 0、returns error",
			rowAffected: 0,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentInValidIDForTest,
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodDELETE,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
		{
			name:        "when RowAffected is 2、returns error",
			rowAffected: 2,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentInValidIDForTest,
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodDELETE,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
		{
			name:        "when DB error has occurred、returns error",
			rowAffected: 0,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.CommentInValidIDForTest,
				err: errors.New(model.ErrorMessageForTest),
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodDELETE,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "DELETE FROM comments WHERE id=\\?"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.id).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.id).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &commentRepository{}

			err := repo.DeleteComment(tt.args.ctx, tt.args.m, tt.args.id)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("commentRepository.DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
