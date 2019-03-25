package application

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	mock_repository "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository/mock"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
	mock_service "github.com/sekky0905/nuxt-vue-go-chat/server/domain/service/mock"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
)

func Test_authenticationService_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		m                 repository.DBManager
		userRepository    repository.UserRepository
		sessionRepository repository.SessionRepository
		userService       service.UserService
		sessionService    service.SessionService
		txCloser          CloseTransaction
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type mockUserRepoArgs struct {
		user *model.User
	}

	type mockSessionRepoArgs struct {
		session *model.Session
	}

	type mockUserServiceArgs struct {
		ctx      context.Context
		name     string
		password string
	}

	type mockSessionServiceArgs struct {
		ctx    context.Context
		id     string
		userID uint32
	}

	type mockUserRepoReturns struct {
		id  uint32
		err error
	}
	type mockUserServiceReturns struct {
		found bool
		err   error
		user  *model.User
	}

	type mockSessionRepoReturns struct {
		err error
	}

	type mockSessionServiceReturns struct {
		found   bool
		err     error
		session *model.Session
	}

	testutil.SetFakeTime(time.Now())

	tests := []struct {
		name   string
		fields fields
		args   args
		mockUserRepoArgs
		mockUserServiceArgs
		mockSessionRepoArgs
		mockSessionServiceArgs
		mockUserRepoReturns
		mockUserServiceReturns
		mockSessionRepoReturns
		mockSessionServiceReturns
		wantUser *model.User
		wantErr  error
	}{
		{
			name: "When appropriate name and password are given and the user which name is same as given name does'nt exist, returns user and nil",
			fields: fields{
				m:                 mock_repository.NewMockDBManager(ctrl),
				userRepository:    mock_repository.NewMockUserRepository(ctrl),
				sessionRepository: mock_repository.NewMockSessionRepository(ctrl),
				userService:       mock_service.NewMockUserService(ctrl),
				sessionService:    mock_service.NewMockSessionService(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			mockUserRepoArgs: mockUserRepoArgs{
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},

			mockSessionRepoArgs: mockSessionRepoArgs{
				session: &model.Session{
					ID:        model.SessionValidIDForTest,
					UserID:    model.UserValidIDForTest,
					CreatedAt: testutil.TimeNow(),
				},
			},
			mockUserServiceArgs: mockUserServiceArgs{
				ctx:      context.Background(),
				name:     model.UserNameForTest,
				password: model.PasswordForTest,
			},
			mockSessionServiceArgs: mockSessionServiceArgs{
				ctx:    context.Background(),
				id:     model.SessionValidIDForTest,
				userID: model.UserValidIDForTest,
			},

			mockUserRepoReturns: mockUserRepoReturns{
				id:  model.UserValidIDForTest,
				err: nil,
			},
			mockSessionRepoReturns: mockSessionRepoReturns{
				err: nil,
			},
			mockUserServiceReturns: mockUserServiceReturns{
				found: false,
				err:   nil,
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockSessionServiceReturns: mockSessionServiceReturns{
				found: false,
				err:   nil,
				session: &model.Session{
					ID:        model.SessionValidIDForTest,
					UserID:    model.UserValidIDForTest,
					CreatedAt: testutil.TimeNow(),
				},
			},
			wantUser: &model.User{
				ID:        model.UserValidIDForTest,
				Name:      model.UserNameForTest,
				SessionID: model.SessionValidIDForTest,
				Password:  model.PasswordForTest,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},

			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, ok := tt.fields.m.(*mock_repository.MockDBManager)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			m.EXPECT().Begin().Return(mock_repository.NewMockTxManager(ctrl), nil)

			us, ok := tt.fields.userService.(*mock_service.MockUserService)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			us.EXPECT().IsAlreadyExistName(tt.args.ctx, tt.mockUserServiceArgs.name).Return(tt.mockUserServiceReturns.found, tt.mockUserServiceReturns.err)
			us.EXPECT().NewUser(tt.mockUserServiceArgs.name, tt.mockUserServiceArgs.password).Return(tt.mockUserServiceReturns.user, tt.mockUserServiceReturns.err)

			ur, ok := tt.fields.userRepository.(*mock_repository.MockUserRepository)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			ur.EXPECT().InsertUser(tt.args.ctx, tt.fields.m, tt.mockUserRepoArgs.user).Return(tt.mockUserRepoReturns.id, tt.mockUserRepoReturns.err)

			ss, ok := tt.fields.sessionService.(*mock_service.MockSessionService)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			ss.EXPECT().IsAlreadyExistID(tt.mockSessionServiceArgs.ctx, tt.mockSessionServiceArgs.id).Return(tt.mockSessionServiceReturns.found, tt.mockSessionServiceReturns.err)
			ss.EXPECT().SessionID().Return(model.SessionValidIDForTest)
			ss.EXPECT().NewSession(tt.mockSessionServiceArgs.userID).Return(tt.mockSessionServiceReturns.session)

			sr, ok := tt.fields.sessionRepository.(*mock_repository.MockSessionRepository)
			if !ok {
				t.Fatal("failed to assert MockSessionRepository")
			}
			sr.EXPECT().InsertSession(tt.args.ctx, tt.fields.m, tt.mockSessionRepoArgs.session).Return(tt.mockSessionRepoReturns.err)

			a := &authenticationService{
				m:                 tt.fields.m,
				userRepository:    tt.fields.userRepository,
				sessionRepository: tt.fields.sessionRepository,
				userService:       tt.fields.userService,
				sessionService:    tt.fields.sessionService,
				txCloser:          tt.fields.txCloser,
			}

			gotUser, err := a.SignUp(tt.args.ctx, tt.args.user)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("authenticationService.SignUp() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("authenticationService.SignUp() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func Test_authenticationService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		m                     repository.DBManager
		userRepository        repository.UserRepository
		sessionRepository     repository.SessionRepository
		userService           service.UserService
		sessionService        service.SessionService
		authenticationService service.AuthenticationService
		txCloser              CloseTransaction
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type mockSessionServiceArgs struct {
		ctx    context.Context
		id     string
		userID uint32
	}

	type mockAuthenticationServiceArgs struct {
		ctx      context.Context
		userName string
		password string
	}

	type mockAuthenticationServiceReturns struct {
		ok   bool
		user *model.User
		err  error
	}

	type mockUserRepoArgs struct {
		user *model.User
	}

	type mockUserRepoReturns struct {
		id  uint32
		err error
	}

	type mockSessionRepoArgs struct {
		ctx     context.Context
		m       repository.DBManager
		session *model.Session
	}

	type mockSessionRepoReturns struct {
		err error
	}

	type mockSessionServiceReturns struct {
		found   bool
		err     error
		session *model.Session
	}

	testutil.SetFakeTime(time.Now())

	tests := []struct {
		name   string
		fields fields
		args   args
		mockUserRepoArgs
		mockUserRepoReturns
		mockSessionRepoArgs
		mockSessionRepoReturns
		mockSessionServiceArgs
		mockSessionServiceReturns
		mockAuthenticationServiceArgs
		mockAuthenticationServiceReturns
		wantUser *model.User
		wantErr  error
	}{
		{
			name: "When appropriate name and password are given and the user which name is same as given name does'nt exist, returns user and nil",
			fields: fields{
				m:                     mock_repository.NewMockDBManager(ctrl),
				userRepository:        mock_repository.NewMockUserRepository(ctrl),
				sessionService:        mock_service.NewMockSessionService(ctrl),
				authenticationService: mock_service.NewMockAuthenticationService(ctrl),
				sessionRepository:     mock_repository.NewMockSessionRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			mockUserRepoArgs: mockUserRepoArgs{
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockUserRepoReturns: mockUserRepoReturns{
				id:  model.UserValidIDForTest,
				err: nil,
			},
			mockSessionRepoArgs: mockSessionRepoArgs{
				ctx: context.Background(),
				m:   mock_repository.NewMockDBManager(ctrl),
				session: &model.Session{
					ID:        model.SessionValidIDForTest,
					UserID:    model.UserValidIDForTest,
					CreatedAt: testutil.TimeNow(),
				},
			},
			mockSessionRepoReturns: mockSessionRepoReturns{
				err: nil,
			},
			mockSessionServiceArgs: mockSessionServiceArgs{
				ctx:    context.Background(),
				id:     model.SessionValidIDForTest,
				userID: model.UserValidIDForTest,
			},
			mockSessionServiceReturns: mockSessionServiceReturns{
				found: false,
				err:   nil,
				session: &model.Session{
					ID:        model.SessionValidIDForTest,
					UserID:    model.UserValidIDForTest,
					CreatedAt: testutil.TimeNow(),
				},
			},
			mockAuthenticationServiceArgs: mockAuthenticationServiceArgs{
				ctx:      context.Background(),
				userName: model.UserNameForTest,
				password: model.PasswordForTest,
			},
			mockAuthenticationServiceReturns: mockAuthenticationServiceReturns{
				ok: true,
				user: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					Password:  model.PasswordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: nil,
			},
			wantUser: &model.User{
				ID:        model.UserValidIDForTest,
				Name:      model.UserNameForTest,
				SessionID: model.SessionValidIDForTest,
				Password:  model.PasswordForTest,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},

			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, ok := tt.fields.m.(*mock_repository.MockDBManager)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			m.EXPECT().Begin().Return(mock_repository.NewMockTxManager(ctrl), nil)

			ur, ok := tt.fields.userRepository.(*mock_repository.MockUserRepository)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			ur.EXPECT().UpdateUser(tt.args.ctx, tt.fields.m, tt.mockUserRepoArgs.user.ID, tt.mockUserRepoArgs.user).Return(tt.mockUserRepoReturns.err)

			sr, ok := tt.fields.sessionRepository.(*mock_repository.MockSessionRepository)
			if !ok {
				t.Fatal("failed to assert MockSessionRepository")
			}
			sr.EXPECT().InsertSession(tt.mockSessionRepoArgs.ctx, tt.mockSessionRepoArgs.m, tt.mockSessionRepoArgs.session).Return(tt.mockSessionRepoReturns.err)

			ss, ok := tt.fields.sessionService.(*mock_service.MockSessionService)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			ss.EXPECT().IsAlreadyExistID(tt.mockSessionServiceArgs.ctx, tt.mockSessionServiceArgs.id).Return(tt.mockSessionServiceReturns.found, tt.mockSessionServiceReturns.err)
			ss.EXPECT().SessionID().Return(model.SessionValidIDForTest)
			ss.EXPECT().NewSession(tt.mockSessionServiceArgs.userID).Return(tt.mockSessionServiceReturns.session)

			as, ok := tt.fields.authenticationService.(*mock_service.MockAuthenticationService)
			if !ok {
				t.Fatal("failed to assert MockSessionRepository")
			}
			as.EXPECT().
				Authenticate(tt.mockAuthenticationServiceArgs.ctx, tt.mockAuthenticationServiceArgs.userName, tt.mockAuthenticationServiceArgs.password).
				Return(tt.mockAuthenticationServiceReturns.ok, tt.mockAuthenticationServiceReturns.user, tt.mockAuthenticationServiceReturns.err)

			a := &authenticationService{
				m:                     tt.fields.m,
				userRepository:        tt.fields.userRepository,
				sessionService:        tt.fields.sessionService,
				authenticationService: tt.fields.authenticationService,
				sessionRepository:     tt.fields.sessionRepository,
				txCloser:              tt.fields.txCloser,
			}

			gotUser, err := a.Login(tt.args.ctx, tt.args.user)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("authenticationService.Login() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("authenticationService.Login() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
