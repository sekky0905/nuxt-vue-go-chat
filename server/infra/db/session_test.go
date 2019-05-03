package db

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Test_sessionRepository_ErrorMsg(t *testing.T) {
	type args struct {
		method model.RepositoryMethod
		err    error
	}

	tests := []struct {
		name    string
		args    args
		wantErr *model.RepositoryError
	}{
		{
			name: "When given appropriate args, returns appropriate error",
			args: args{
				method: model.RepositoryMethodInsert,
				err:    errors.New(model.ErrorMessageForTest),
			},
			wantErr: &model.RepositoryError{
				BaseErr:          errors.New(model.ErrorMessageForTest),
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameSession,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &sessionRepository{}
			if err := repo.ErrorMsg(tt.args.method, tt.args.err); errors.Cause(err).Error() != tt.wantErr.Error() {
				t.Errorf("sessionRepository.ErrorMsg() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}

func Test_sessionRepository_GetSessionByID(t *testing.T) {
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
		id  string
	}

	tests := []struct {
		name    string
		args    args
		want    *model.Session
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When a session specified by id exists, returns a session",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.SessionValidIDForTest,
			},
			want: &model.Session{
				ID:        model.SessionValidIDForTest,
				UserID:    model.UserValidIDForTest,
				CreatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "When a session specified by id does not exist, returns NoSuchDataError",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.SessionInValidIDForTest,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyName:    model.IDProperty,
				PropertyValue:   model.SessionInValidIDForTest,
				DomainModelName: model.DomainModelNameSession,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := "SELECT id, user_id, created_at FROM sessions WHERE id=?"
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "user_id", "created_at"}).
					AddRow(tt.want.ID, tt.want.UserID, tt.want.CreatedAt)
				prep.ExpectQuery().WithArgs(tt.want.ID).WillReturnRows(rows)
			}

			repo := &sessionRepository{}
			got, err := repo.GetSessionByID(tt.args.ctx, tt.args.m, tt.args.id)

			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("sessionRepository.GetSessionByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sessionRepository.GetSessionByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sessionRepository_InsertSession(t *testing.T) {
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
		session *model.Session
		err     error
	}

	tests := []struct {
		name        string
		args        args
		rowAffected int64
		want        string
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a session which has ID, User_ID, CreatedAt is given, returns ID",
			args: args{
				ctx: context.Background(),
				m:   db,
				session: &model.Session{
					ID:        model.SessionValidIDForTest,
					UserID:    model.UserValidIDForTest,
					CreatedAt: testutil.TimeNow(),
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
				session: &model.Session{
					ID:        model.SessionInValidIDForTest,
					UserID:    model.UserValidIDForTest,
					CreatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameSession,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				session: &model.Session{
					ID:        model.SessionInValidIDForTest,
					UserID:    model.UserValidIDForTest,
					CreatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameSession,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				session: &model.Session{
					ID:        model.SessionInValidIDForTest,
					UserID:    model.UserValidIDForTest,
					CreatedAt: testutil.TimeNow(),
				},
				err: errors.New(model.ErrorMessageForTest),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameSession,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "INSERT INTO sessions"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.session.ID, tt.args.session.UserID, tt.args.session.CreatedAt).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.session.ID, tt.args.session.UserID, tt.args.session.CreatedAt).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &sessionRepository{}

			err := repo.InsertSession(tt.args.ctx, tt.args.m, tt.args.session)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("sessionRepository.InsertSession() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func Test_sessionRepository_DeleteSession(t *testing.T) {
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
			name:        "When a user specified by id exists, returns nil",
			rowAffected: 1,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.UserValidIDForTest,
			},
			wantErr: nil,
		},
		{
			name:        "when RowAffected is 0、returns error",
			rowAffected: 0,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.UserInValidIDForTest,
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodDELETE,
				DomainModelName:  model.DomainModelNameUser,
			},
		},
		{
			name:        "when RowAffected is 2、returns error",
			rowAffected: 2,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.UserInValidIDForTest,
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodDELETE,
				DomainModelName:  model.DomainModelNameUser,
			},
		},
		{
			name:        "when DB error has occurred、returns error",
			rowAffected: 0,
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.UserInValidIDForTest,
				err: errors.New(model.ErrorMessageForTest),
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodDELETE,
				DomainModelName:  model.DomainModelNameUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "DELETE FROM users WHERE id=\\?"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.id).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.id).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &userRepository{}

			err := repo.DeleteUser(tt.args.ctx, tt.args.m, tt.args.id)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("userRepository.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
