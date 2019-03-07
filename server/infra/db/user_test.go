package db

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewUserRepository(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want repository.UserRepository
	}{
		{
			name: "When given appropriate args, returns UserRepository",
			args: args{
				ctx: context.Background(),
			},
			want: &userRepository{
				context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepository(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_ErrorMsg(t *testing.T) {
	const errMsg = "test"

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
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepository{
				ctx: tt.fields.ctx,
			}
			if err := repo.ErrorMsg(tt.args.method, tt.args.err); errors.Cause(err).Error() != tt.wantErr.Error() {
				t.Errorf("userRepository.ErrorMsg() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepository_GetUserByID(t *testing.T) {
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
		id uint32
	}

	var validID uint32 = 1
	var inValidID uint32 = 2

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When specified user exists, returns a user",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: validID,
			},
			want: &model.User{
				ID:        validID,
				Name:      "test",
				SessionID: "test12345678",
				Password:  "test",
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "When specified user does not exist, returns NoSuchDataError",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: inValidID,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
				PropertyNameForUser:         model.IDPropertyForUser,
				PropertyValue:               inValidID,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := "SELECT id, name, session_id, password, created_at, updated_at FROM users WHERE id=?"
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name", "session_id", "password", "created_at", "updated_at"}).
					AddRow(tt.want.ID, tt.want.Name, tt.want.SessionID, tt.want.Password, tt.want.CreatedAt, tt.want.UpdatedAt)
				prep.ExpectQuery().WithArgs(tt.want.ID).WillReturnRows(rows)
			}

			repo := &userRepository{
				ctx: tt.fields.ctx,
			}
			got, err := repo.GetUserByID(tt.args.m, tt.args.id)

			if tt.wantErr != nil {
				if !reflect.DeepEqual(err, tt.wantErr) {
					t.Errorf("userRepository.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepository.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
