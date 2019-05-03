package db

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewThreadRepository(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want repository.ThreadRepository
	}{
		{
			name: "When given appropriate args, returns ThreadRepository",
			args: args{
				ctx: context.Background(),
			},
			want: &threadRepository{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewThreadRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewThreadRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_threadRepository_ErrorMsg(t *testing.T) {
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
				BaseErr:          errors.New(model.ErrorMessageForTest),
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &threadRepository{}
			if err := repo.ErrorMsg(tt.args.method, tt.args.err); errors.Cause(err).Error() != tt.wantErr.Error() {
				t.Errorf("threadRepository.ErrorMsg() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}

func Test_threadRepository_ListThreads(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		return
	}

	defer db.Close()

	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx    context.Context
		m      query.SQLManager
		limit  int
		cursor uint32
	}

	tests := []struct {
		name       string
		repo       *threadRepository
		args       args
		want       *model.ThreadList
		returnMock []*model.Thread
		wantErr    error
	}{
		{
			name: "When limit = 20, cursor = 1 are given and there are over 21 data, ListThreads returns ThreadList which has Threads(ID: 1~20), HasNext = yes, Cursor = 21",
			repo: &threadRepository{},
			args: args{
				ctx:    context.Background(),
				m:      db,
				limit:  20,
				cursor: 1,
			},
			want: &model.ThreadList{
				Threads: testutil.GenerateThreadHelper(1, 20),
				HasNext: true,
				Cursor:  21,
			},
			returnMock: testutil.GenerateThreadHelper(1, 21),
			wantErr:    nil,
		},
		{
			name: "When limit = 20, cursor = 21 are given and there are over 41 data, ListThreads returns ThreadList which has 21 Threads(ID: 21~40), HasNext = yes, Cursor = 41",
			repo: &threadRepository{},
			args: args{
				ctx:    context.Background(),
				m:      db,
				limit:  20,
				cursor: 21,
			},
			want: &model.ThreadList{
				Threads: testutil.GenerateThreadHelper(21, 40),
				HasNext: true,
				Cursor:  41,
			},
			returnMock: testutil.GenerateThreadHelper(21, 41),
			wantErr:    nil,
		},
		{
			name: "When limit = 20, cursor = 1 are given and there are over 10 data, ListThreads returns ThreadList which has 10 Threads(ID: 1~10), HasNext = yes, Cursor = 10",
			repo: &threadRepository{},
			args: args{
				ctx:    context.Background(),
				m:      db,
				limit:  20,
				cursor: 1,
			},
			want: &model.ThreadList{
				Threads: testutil.GenerateThreadHelper(1, 10),
				HasNext: false,
				Cursor:  0,
			},
			returnMock: testutil.GenerateThreadHelper(1, 10),
			wantErr:    nil,
		},
		{
			name: "When limit = 20, cursor = 1 are given and there are no data, ListThreads returns error",
			repo: &threadRepository{},
			args: args{
				ctx:    context.Background(),
				m:      db,
				limit:  20,
				cursor: 1,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				BaseErr:         err,
				DomainModelName: model.DomainModelNameThread,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `SELECT (.+)
	FROM threads AS t
	INNER JOIN users AS u
	(.+);`
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WithArgs(tt.args.cursor, readyLimitForHasNext(tt.args.limit)).WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"t.id", "t.title", "u.id", "u.name", "t.created_at", "t.updated_at"})

				for _, thread := range tt.returnMock {
					rows.AddRow(thread.ID, thread.Title, thread.User.ID, thread.User.Name, thread.CreatedAt, thread.UpdatedAt)
				}

				prep.ExpectQuery().WithArgs(tt.args.cursor, readyLimitForHasNext(tt.args.limit)).WillReturnRows(rows)
			}

			repo := &threadRepository{}
			got, err := repo.ListThreads(tt.args.ctx, tt.args.m, tt.args.cursor, tt.args.limit)
			if tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("threadRepository.ListThreads() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("threadRepository.ListThreads() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_threadRepository_GetThreadByID(t *testing.T) {
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
		want    *model.Thread
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When a thread specified by id exists, returns a thread",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.ThreadValidIDForTest,
			},
			want: &model.Thread{
				ID:    model.ThreadValidIDForTest,
				Title: model.TitleForTest,
				User: &model.User{
					ID:   model.UserValidIDForTest,
					Name: model.UserNameForTest,
				},
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "When a thread specified by id does not exist, returns NoSuchDataError",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.ThreadInValidIDForTest,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyName:    model.IDProperty,
				PropertyValue:   model.UserInValidIDForTest,
				DomainModelName: model.DomainModelNameThread,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `SELECT (.+)
	FROM threads AS t
	INNER JOIN users AS u
	(.+);`
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"t.id", "t.title", "u.id", "u.name", "t.created_at", "t.updated_at"}).
					AddRow(tt.want.ID, tt.want.Title, tt.want.User.ID, tt.want.User.Name, tt.want.CreatedAt, tt.want.UpdatedAt)
				prep.ExpectQuery().WithArgs(tt.want.ID).WillReturnRows(rows)
			}

			repo := &threadRepository{}
			got, err := repo.GetThreadByID(tt.args.ctx, tt.args.m, tt.args.id)

			if tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("threadRepository.GetThreadByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("threadRepository.GetThreadByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_threadRepository_GetThreadByTitle(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		return
	}

	defer db.Close()

	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx   context.Context
		m     query.SQLManager
		title string
	}

	tests := []struct {
		name    string
		args    args
		want    *model.Thread
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When a thread specified by id exists, returns a thread",
			args: args{
				ctx:   context.Background(),
				m:     db,
				title: model.TitleForTest,
			},
			want: &model.Thread{
				ID:    model.ThreadValidIDForTest,
				Title: model.TitleForTest,
				User: &model.User{
					ID:   model.UserValidIDForTest,
					Name: model.UserNameForTest,
				},
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "When a thread specified by id does not exist, returns NoSuchDataError",
			args: args{
				ctx:   context.Background(),
				m:     db,
				title: model.TitleForTest,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyName:    model.NameProperty,
				PropertyValue:   model.TitleForTest,
				DomainModelName: model.DomainModelNameThread,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `SELECT (.+)
	FROM threads AS t
	INNER JOIN users AS u
	(.+);`
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"t.id", "t.title", "u.id", "u.name", "t.created_at", "t.updated_at"}).
					AddRow(tt.want.ID, tt.want.Title, tt.want.User.ID, tt.want.User.Name, tt.want.CreatedAt, tt.want.UpdatedAt)
				prep.ExpectQuery().WithArgs(tt.want.Title).WillReturnRows(rows)
			}

			repo := &threadRepository{}
			got, err := repo.GetThreadByTitle(tt.args.ctx, tt.args.m, tt.args.title)

			if tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("threadRepository.GetThreadByTitle() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("threadRepository.GetThreadByTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_threadRepository_InsertThread(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx    context.Context
		m      query.SQLManager
		thread *model.Thread
		err    error
	}

	tests := []struct {
		name        string
		args        args
		rowAffected int64
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a thread which has ID, Name, Title, User, CreatedAt, UpdatedAt is given, returns ID",
			args: args{
				ctx: context.Background(),
				m:   db,
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
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
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: errors.New(model.ErrorMessageForTest),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "INSERT INTO threads"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.thread.ID, tt.args.thread.Name, tt.args.thread.SessionID, tt.args.thread.Password, tt.args.thread.CreatedAt, tt.args.thread.UpdatedAt).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.thread.ID, tt.args.thread.Name, tt.args.thread.SessionID, tt.args.thread.Password, tt.args.thread.CreatedAt, tt.args.thread.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &threadRepository{}

			_, err := repo.InsertThread(tt.args.ctx, tt.args.m, tt.args.thread)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("threadRepository.InsertThread() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func Test_threadRepository_UpdateThread(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx    context.Context
		m      query.SQLManager
		id     uint32
		thread *model.Thread
		err    error
	}

	tests := []struct {
		name        string
		args        args
		rowAffected int64
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a thread which has ID, Name, Title, User, CreatedAt, UpdatedAt is given, returns nil",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.ThreadValidIDForTest,
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
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
				id:  model.ThreadInValidIDForTest,
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodUPDATE,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.ThreadInValidIDForTest,
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodUPDATE,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.ThreadInValidIDForTest,
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: errors.New(model.ErrorMessageForTest),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodUPDATE,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "UPDATE threads SET title=\\?, updated_at=\\? WHERE id=\\?"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.thread.Title, tt.args.thread.UpdatedAt, tt.args.id).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.thread.Title, tt.args.thread.UpdatedAt, tt.args.id).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &threadRepository{}
			err := repo.UpdateThread(tt.args.ctx, tt.args.m, tt.args.id, tt.args.thread)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("threadRepository.UpdateThread() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func Test_threadRepository_DeleteThread(t *testing.T) {
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
			name:        "When a thread specified by id exists, returns nil",
			rowAffected: 1,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.ThreadValidIDForTest,
			},
			wantErr: nil,
		},
		{
			name:        "when RowAffected is 0、returns error",
			rowAffected: 0,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.ThreadInValidIDForTest,
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodDELETE,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
		{
			name:        "when RowAffected is 2、returns error",
			rowAffected: 2,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.ThreadInValidIDForTest,
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodDELETE,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
		{
			name:        "when DB error has occurred、returns error",
			rowAffected: 0,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.ThreadInValidIDForTest,
				err: errors.New(model.ErrorMessageForTest),
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodDELETE,
				DomainModelName:  model.DomainModelNameThread,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "DELETE FROM threads WHERE id=\\?"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.id).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.id).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &threadRepository{}

			err := repo.DeleteThread(tt.args.ctx, tt.args.m, tt.args.id)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("threadRepository.DeleteThread() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
