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
			want: &userRepository{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_ErrorMsg(t *testing.T) {
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
				DomainModelName:  model.DomainModelNameUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepository{}
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

	type args struct {
		ctx context.Context
		m   query.SQLManager
		id  uint32
	}

	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When a user specified by id exists, returns a user",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.UserValidIDForTest,
			},
			want: &model.User{
				ID:        model.UserValidIDForTest,
				Name:      model.UserNameForTest,
				SessionID: model.SessionValidIDForTest,
				Password:  model.PasswordForTest,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "When a user specified by id does not exist, returns NoSuchDataError",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.UserInValidIDForTest,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyName:    model.IDProperty,
				PropertyValue:   model.UserInValidIDForTest,
				DomainModelName: model.DomainModelNameUser,
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

			repo := &userRepository{}
			got, err := repo.GetUserByID(tt.args.ctx, tt.args.m, tt.args.id)

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

func Test_userRepository_GetUserByName(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		return
	}

	defer db.Close()

	type args struct {
		ctx  context.Context
		m    query.SQLManager
		name string
	}

	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When a user specified by name exists, returns a user",
			args: args{
				ctx:  context.Background(),
				m:    db,
				name: model.UserNameForTest,
			},
			want: &model.User{
				Name:      model.UserNameForTest,
				SessionID: model.SessionValidIDForTest,
				Password:  model.PasswordForTest,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "When a user specified by name does not exist, returns NoSuchDataError",
			args: args{
				ctx:  context.Background(),
				m:    db,
				name: "test2",
			},
			want: nil,

			wantErr: &model.NoSuchDataError{
				PropertyName:    model.NameProperty,
				PropertyValue:   "test2",
				DomainModelName: model.DomainModelNameUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := "SELECT id, name, session_id, password, created_at, updated_at FROM users WHERE name=?"
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name", "session_id", "password", "created_at", "updated_at"}).
					AddRow(tt.want.ID, tt.want.Name, tt.want.SessionID, tt.want.Password, tt.want.CreatedAt, tt.want.UpdatedAt)
				prep.ExpectQuery().WithArgs(tt.want.Name).WillReturnRows(rows)
			}

			repo := &userRepository{}
			got, err := repo.GetUserByName(tt.args.ctx, tt.args.m, tt.args.name)
			if tt.wantErr != nil {
				if !reflect.DeepEqual(err, tt.wantErr) {
					t.Errorf("userRepository.GetUserByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepository.GetUserByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_InsertUser(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx  context.Context
		m    query.SQLManager
		user *model.User
		err  error
	}

	tests := []struct {
		name        string
		args        args
		rowAffected int64
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a user which has ID, Name, Session_ID, Password, CreatedAt, UpdatedAt is given, returns ID",
			args: args{
				ctx: context.Background(),
				m:   db,
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
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
				user: &model.User{
					ID:        model.UserInValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameUser,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				user: &model.User{
					ID:        model.UserInValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameUser,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				user: &model.User{
					ID:        model.UserInValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: errors.New(model.ErrorMessageForTest),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodInsert,
				DomainModelName:  model.DomainModelNameUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "INSERT INTO users"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.user.ID, tt.args.user.Name, tt.args.user.SessionID, tt.args.user.Password, tt.args.user.CreatedAt, tt.args.user.UpdatedAt).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.user.ID, tt.args.user.Name, tt.args.user.SessionID, tt.args.user.Password, tt.args.user.CreatedAt, tt.args.user.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &userRepository{}

			_, err := repo.InsertUser(tt.args.ctx, tt.args.m, tt.args.user)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("userRepository.InsertUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func Test_userRepository_UpdateUser(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	testutil.SetFakeTime(time.Now())

	type args struct {
		ctx  context.Context
		m    query.SQLManager
		id   uint32
		user *model.User
		err  error
	}

	tests := []struct {
		name        string
		args        args
		rowAffected int64
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a user which has Name, Session_ID, Password, UpdatedAt is given, returns nil",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.UserValidIDForTest,
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
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
				id:  model.UserInValidIDForTest,
				user: &model.User{
					ID:        model.UserInValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodUPDATE,
				DomainModelName:  model.DomainModelNameUser,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.UserInValidIDForTest,
				user: &model.User{
					ID:        model.UserInValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodUPDATE,
				DomainModelName:  model.DomainModelNameUser,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			args: args{
				ctx: context.Background(),
				m:   db,
				id:  model.UserInValidIDForTest,
				user: &model.User{
					ID:        model.UserInValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: errors.New(model.ErrorMessageForTest),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod: model.RepositoryMethodUPDATE,
				DomainModelName:  model.DomainModelNameUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "UPDATE users SET session_id=\\?, password=\\?, created_at=\\?, updated_at=\\? WHERE id=\\?"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.user.SessionID, tt.args.user.Password, tt.args.user.CreatedAt, tt.args.user.UpdatedAt, tt.args.id).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.user.SessionID, tt.args.user.Password, tt.args.user.CreatedAt, tt.args.user.UpdatedAt, tt.args.id).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &userRepository{}
			err := repo.UpdateUser(tt.args.ctx, tt.args.m, tt.args.id, tt.args.user)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("userRepository.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func Test_userRepository_DeleteUser(t *testing.T) {
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
