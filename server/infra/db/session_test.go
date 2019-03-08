package db

import (
	"context"
	"reflect"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
)

const (
	sessionValidIDForTest   = "testValidSessionID12345678"
	sessionInValidIDForTest = "testInvalidSessionID12345678"
)

func TestNewSessionRepository(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want repository.SessionRepository
	}{
		{
			name: "When given appropriate args, returns SessionRepository",
			args: args{
				ctx: context.Background(),
			},
			want: &sessionRepository{
				context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSessionRepository(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSessionRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sessionRepository_ErrorMsg(t *testing.T) {
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
				err:    errors.New(errMsg),
			},
			wantErr: &model.RepositoryError{
				BaseErr:                     errors.New(errMsg),
				RepositoryMethod:            model.RepositoryMethodInsert,
				DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameSessionForUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &sessionRepository{
				ctx: tt.fields.ctx,
			}
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

	type fields struct {
		ctx context.Context
	}
	type args struct {
		m  repository.DBManager
		id string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Session
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When a session specified by id exists, returns a session",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: sessionValidIDForTest,
			},
			want: &model.Session{
				ID:        sessionValidIDForTest,
				UserID:    userValidIDForTest,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "When a session specified by id does not exist, returns NoSuchDataError",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: sessionInValidIDForTest,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
				PropertyNameForUser:         model.IDPropertyForUser,
				PropertyValue:               sessionInValidIDForTest,
				DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameSessionForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := "SELECT id, user_id, created_at, updated_at FROM sessions WHERE id=?"
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).
					AddRow(tt.want.ID, tt.want.UserID, tt.want.CreatedAt, tt.want.UpdatedAt)
				prep.ExpectQuery().WithArgs(tt.want.ID).WillReturnRows(rows)
			}

			repo := &sessionRepository{
				ctx: tt.fields.ctx,
			}
			got, err := repo.GetSessionByID(tt.args.m, tt.args.id)

			if tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
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
