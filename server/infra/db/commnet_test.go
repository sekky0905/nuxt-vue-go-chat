package db

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	. "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
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
		ctx      context.Context
		m        SQLManager
		threadID uint32
		limit    int
		cursor   uint32
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
				ctx:      context.Background(),
				m:        db,
				threadID: model.ThreadValidIDForTest,
				limit:    20,
				cursor:   1,
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
				ctx:      context.Background(),
				m:        db,
				threadID: model.ThreadValidIDForTest,
				limit:    20,
				cursor:   21,
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
				ctx:      context.Background(),
				m:        db,
				threadID: model.ThreadValidIDForTest,
				limit:    20,
				cursor:   1,
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
				ctx:      context.Background(),
				m:        db,
				threadID: model.ThreadValidIDForTest,
				limit:    20,
				cursor:   1,
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
				rows := sqlmock.NewRows([]string{"c.id", "c.content", "u.id", "u.name", "c.thread_id", "c.created_at", "c.updated_at"})

				for _, comment := range tt.returnMock {
					rows.AddRow(comment.ID, comment.Content, comment.User.ID, comment.User.Name, comment.ThreadID, comment.CreatedAt, comment.UpdatedAt)
				}

				prep.ExpectQuery().WithArgs(tt.args.cursor, tt.args.threadID, readyLimitForHasNext(tt.args.limit)).WillReturnRows(rows)
			}

			repo := &commentRepository{}
			got, err := repo.ListComments(tt.args.ctx, tt.args.m, tt.args.threadID, tt.args.limit, tt.args.cursor)
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
