package service

import (
	"context"
	"testing"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	mock_repository "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository/mock"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
)

func Test_userService_IsAlreadyExistID(t *testing.T) {
	// for gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock_repository.NewMockUserRepository(ctrl)

	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx context.Context
		m   query.DBManager
		id  uint32
	}

	type returnArgs struct {
		user *model.User
		err  error
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		returnArgs
		want    bool
		wantErr error
	}{
		{
			name: "When specified user already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx: context.Background(),
				m:   db.NewDBManager(),
				id:  model.UserValidIDForTest,
			},
			returnArgs: returnArgs{
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: nil,
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "When specified user doesn't already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx: context.Background(),
				m:   db.NewDBManager(),
				id:  model.UserInValidIDForTest,
			},
			returnArgs: returnArgs{
				user: nil,
				err:  nil,
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "When some error has occurred, return false and error.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx: context.Background(),
				m:   db.NewDBManager(),
				id:  model.UserInValidIDForTest,
			},
			returnArgs: returnArgs{
				user: nil,
				err:  errors.New(model.ErrorMessageForTest),
			},
			want:    false,
			wantErr: errors.New(model.ErrorMessageForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				repo: mock,
			}

			mock.EXPECT().GetUserByID(tt.args.ctx, tt.args.m, tt.args.id).Return(tt.returnArgs.user, tt.returnArgs.err)

			got, err := s.IsAlreadyExistID(tt.args.ctx, tt.args.m, tt.args.id)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("userService.IsAlreadyExistID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("userService.IsAlreadyExistID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_IsAlreadyExistName(t *testing.T) {
	// for gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock_repository.NewMockUserRepository(ctrl)

	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx  context.Context
		m    query.DBManager
		name string
	}

	type returnArgs struct {
		user *model.User
		err  error
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		returnArgs
		want    bool
		wantErr error
	}{
		{
			name: "When specified user already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx:  context.Background(),
				name: model.UserNameForTest,
			},
			returnArgs: returnArgs{
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: nil,
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "When specified user doesn't already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx:  context.Background(),
				m:    db.NewDBManager(),
				name: model.UserNameForTest,
			},
			returnArgs: returnArgs{
				user: nil,
				err:  nil,
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "When some error has occurred, return false and error.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx:  context.Background(),
				m:    db.NewDBManager(),
				name: model.UserNameForTest,
			},
			returnArgs: returnArgs{
				user: nil,
				err:  errors.New(model.ErrorMessageForTest),
			},
			want:    false,
			wantErr: errors.New(model.ErrorMessageForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				repo: mock,
			}

			mock.EXPECT().GetUserByName(tt.args.ctx, tt.args.m, tt.args.name).Return(tt.returnArgs.user, tt.returnArgs.err)

			got, err := s.IsAlreadyExistName(tt.args.ctx, tt.args.m, tt.args.name)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("userService.IsAlreadyExistName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("userService.IsAlreadyExistName() = %v, want %v", got, tt.want)
			}
		})
	}
}
