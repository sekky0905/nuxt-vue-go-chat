package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/sekky0905/nuxt-vue-go-chat/server/util"

	"github.com/golang/mock/gomock"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	mock_repository "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository/mock"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
)

func Test_authenticationService_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testutil.SetFakeTime(time.Now())

	mockDBM := db.NewDBManager()
	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	type args struct {
		ctx      context.Context
		userName string
		password string
	}

	type mockReturns struct {
		user *model.User
		err  error
	}

	hashedPass, err := util.HashPassword(model.PasswordForTest)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		args     args
		wantOk   bool
		wantUser *model.User
		wantErr  error
		mockReturns
	}{
		{
			name: "When appropriate name and password are given, returns false and nil",
			args: args{
				ctx:      context.Background(),
				userName: model.UserNameForTest,
				password: model.PasswordForTest,
			},
			wantOk: true,
			wantUser: &model.User{
				ID:        model.UserValidIDForTest,
				Name:      model.UserNameForTest,
				SessionID: model.SessionValidIDForTest,
				Password:  hashedPass,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
			mockReturns: mockReturns{
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  hashedPass,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: nil,
			},
		},
		{
			name: "When inappropriate password is given, returns false and nil",
			args: args{
				ctx:      context.Background(),
				userName: model.UserNameForTest,
				password: "invalidPassword",
			},
			wantOk:   false,
			wantUser: nil,
			wantErr:  nil,
			mockReturns: mockReturns{
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  hashedPass,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: nil,
			},
		},
		{
			name: "When user doesn't exist, returns false and error",
			args: args{
				ctx:      context.Background(),
				userName: model.UserNameForTest,
				password: model.PasswordForTest,
			},
			wantOk:   false,
			wantUser: nil,
			wantErr: &model.AuthenticationErr{
				BaseErr: &model.NoSuchDataError{
					BaseErr:         nil,
					PropertyName:    model.NameProperty,
					PropertyValue:   "test",
					DomainModelName: model.DomainModelNameUser,
				},
			},
			mockReturns: mockReturns{
				user: nil,
				err: &model.NoSuchDataError{
					BaseErr:         nil,
					PropertyName:    model.NameProperty,
					PropertyValue:   "test",
					DomainModelName: model.DomainModelNameUser,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authenticationService{
				repo: mockRepo,
			}

			ur, ok := s.repo.(*mock_repository.MockUserRepository)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}

			ur.EXPECT().GetUserByName(tt.args.ctx, mockDBM, tt.args.userName).Return(tt.mockReturns.user, tt.mockReturns.err)

			gotOk, gotUser, err := s.Authenticate(tt.args.ctx, mockDBM, tt.args.userName, tt.args.password)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("authenticationService.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("authenticationService.Authenticate() = %v, want %v", gotOk, tt.wantOk)
			}

			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("authenticationService.Authenticate() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
