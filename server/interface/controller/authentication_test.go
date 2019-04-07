package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	mock_application "github.com/sekky0905/nuxt-vue-go-chat/server/application/mock"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
)

func Test_authenticationController_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		aApp application.AuthenticationService
	}

	type args struct {
		user *model.User
	}

	type errBody struct {
		errCode ErrCode
		filed   string
	}

	type want struct {
		statusCode int
		cookie     string
		body       *UserDTO
		errBody
	}

	type mockArgs struct {
		ctx  context.Context
		user *model.User
	}

	type mockReturns struct {
		user *model.User
		err  error
	}

	testutil.SetFakeTime(time.Date(2019, time.March, 31, 0, 0, 0, 0, time.UTC))

	tests := []struct {
		name   string
		fields fields
		args
		want
		mockArgs
		mockReturns
	}{
		{
			name: "",
			fields: fields{
				aApp: mock_application.NewMockAuthenticationService(ctrl),
			},
			args: args{
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			want: want{
				statusCode: http.StatusOK,
				cookie:     "SESSION_ID=testValidSessionID12345678; Path=/; Max-Age=86400; HttpOnly",
				body: &UserDTO{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				errBody: errBody{},
			},
			mockArgs: mockArgs{
				ctx: context.Background(),
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			mockReturns: mockReturns{
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
		},
		{
			name: "When name is empty, SingUp returns status code 400",
			fields: fields{
				aApp: mock_application.NewMockAuthenticationService(ctrl),
			},
			args: args{
				user: &model.User{
					Name:     "",
					Password: model.PasswordForTest,
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				cookie:     "",
				body:       nil,
				errBody: errBody{
					errCode: InvalidParameterValueFailure,
					filed:   model.NamePropertyForDeveloper.String(),
				},
			},
			mockArgs: mockArgs{
				ctx: context.Background(),
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			mockReturns: mockReturns{
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
		},
		{
			name: "When password is empty, SingUp returns status code 400",
			fields: fields{
				aApp: mock_application.NewMockAuthenticationService(ctrl),
			},
			args: args{
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: "",
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				cookie:     "",
				body:       nil,
				errBody: errBody{
					errCode: InvalidParameterValueFailure,
					filed:   model.PassWordPropertyForDeveloper.String(),
				},
			},
			mockArgs: mockArgs{
				ctx: context.Background(),
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			mockReturns: mockReturns{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aa, ok := tt.fields.aApp.(*mock_application.MockAuthenticationService)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}

			if tt.want.statusCode != http.StatusBadRequest {
				aa.EXPECT().SignUp(tt.mockArgs.ctx, tt.mockArgs.user).Return(tt.mockReturns.user, tt.mockReturns.err)
			}

			ac := NewAuthenticationController(tt.fields.aApp)
			r := gin.New()

			r.POST("/singUp", ac.SignUp)

			rec := httptest.NewRecorder()
			b, err := json.Marshal(tt.args.user)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/singUp", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}
			r.ServeHTTP(rec, req)

			if rec.Code != tt.want.statusCode {
				t.Errorf("status code = %v, want %v", rec.Code, tt.want.statusCode)
				return
			}

			gotCookieVal := rec.Header().Get("Set-Cookie")
			if gotCookieVal != tt.want.cookie {
				t.Errorf("cookie = %v, want %v", gotCookieVal, tt.want.cookie)
				return
			}

			if tt.want.errBody.filed == "" {
				bBody := rec.Body.Bytes()
				uDTO := &UserDTO{}
				if err := json.Unmarshal(bBody, uDTO); err != nil {
					t.Fatal(err)
					return
				}

				if !reflect.DeepEqual(uDTO, tt.want.body) {
					t.Errorf("body = %#v, want %#v", uDTO, tt.want.body)
					return
				}
			} else {
				sBody := rec.Body.String()
				if !strings.Contains(sBody, string(tt.want.errBody.errCode)) {
					t.Errorf("body = %#v, want %#v", sBody, tt.want.errBody.errCode)
					return
				}

				if !strings.Contains(sBody, tt.want.errBody.filed) {
					t.Errorf("body = %#v, want %#v", sBody, tt.want.errBody.filed)
					return
				}
			}
		})
	}
}

func Test_authenticationController_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		aApp application.AuthenticationService
	}

	type args struct {
		user *model.User
	}

	type errBody struct {
		errCode ErrCode
		filed   string
	}

	type want struct {
		statusCode int
		cookie     string
		body       *UserDTO
		errBody
	}

	type mockArgs struct {
		ctx  context.Context
		user *model.User
	}

	type mockReturns struct {
		user *model.User
		err  error
	}

	testutil.SetFakeTime(time.Date(2019, time.March, 31, 0, 0, 0, 0, time.UTC))

	tests := []struct {
		name   string
		fields fields
		args
		want
		mockArgs
		mockReturns
	}{
		{
			name: "",
			fields: fields{
				aApp: mock_application.NewMockAuthenticationService(ctrl),
			},
			args: args{
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			want: want{
				statusCode: http.StatusOK,
				cookie:     "SESSION_ID=testValidSessionID12345678; Path=/; Max-Age=86400; HttpOnly",
				body: &UserDTO{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					SessionID: model.SessionValidIDForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				errBody: errBody{},
			},
			mockArgs: mockArgs{
				ctx: context.Background(),
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			mockReturns: mockReturns{
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
		},
		{
			name: "When name is empty, Login returns status code 400",
			fields: fields{
				aApp: mock_application.NewMockAuthenticationService(ctrl),
			},
			args: args{
				user: &model.User{
					Name:     "",
					Password: model.PasswordForTest,
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				cookie:     "",
				body:       nil,
				errBody: errBody{
					errCode: InvalidParameterValueFailure,
					filed:   model.NamePropertyForDeveloper.String(),
				},
			},
			mockArgs: mockArgs{
				ctx: context.Background(),
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			mockReturns: mockReturns{
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
		},
		{
			name: "When password is empty, Login returns status code 400",
			fields: fields{
				aApp: mock_application.NewMockAuthenticationService(ctrl),
			},
			args: args{
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: "",
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				cookie:     "",
				body:       nil,
				errBody: errBody{
					errCode: InvalidParameterValueFailure,
					filed:   model.PassWordPropertyForDeveloper.String(),
				},
			},
			mockArgs: mockArgs{
				ctx: context.Background(),
				user: &model.User{
					Name:     model.UserNameForTest,
					Password: model.PasswordForTest,
				},
			},
			mockReturns: mockReturns{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aa, ok := tt.fields.aApp.(*mock_application.MockAuthenticationService)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}

			if tt.want.statusCode != http.StatusBadRequest {
				aa.EXPECT().Login(tt.mockArgs.ctx, tt.mockArgs.user).Return(tt.mockReturns.user, tt.mockReturns.err)
			}

			ac := NewAuthenticationController(tt.fields.aApp)
			r := gin.New()

			r.POST("/login", ac.Login)

			rec := httptest.NewRecorder()
			b, err := json.Marshal(tt.args.user)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}
			r.ServeHTTP(rec, req)

			if rec.Code != tt.want.statusCode {
				t.Errorf("status code = %v, want %v", rec.Code, tt.want.statusCode)
				return
			}

			gotCookieVal := rec.Header().Get("Set-Cookie")
			if gotCookieVal != tt.want.cookie {
				t.Errorf("cookie = %v, want %v", gotCookieVal, tt.want.cookie)
				return
			}

			if tt.want.errBody.filed == "" {
				bBody := rec.Body.Bytes()
				uDTO := &UserDTO{}
				if err := json.Unmarshal(bBody, uDTO); err != nil {
					t.Fatal(err)
					return
				}

				if !reflect.DeepEqual(uDTO, tt.want.body) {
					t.Errorf("body = %#v, want %#v", uDTO, tt.want.body)
					return
				}
			} else {
				sBody := rec.Body.String()
				if !strings.Contains(sBody, string(tt.want.errBody.errCode)) {
					t.Errorf("body = %#v, want %#v", sBody, tt.want.errBody.errCode)
					return
				}

				if !strings.Contains(sBody, tt.want.errBody.filed) {
					t.Errorf("body = %#v, want %#v", sBody, tt.want.errBody.filed)
					return
				}
			}
		})
	}
}

func Test_authenticationController_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type errBody struct {
		errCode ErrCode
		filed   string
	}

	type want struct {
		statusCode int
		cookie     string
		errBody
	}

	type mockArgs struct {
		ctx       context.Context
		sessionID string
	}

	type mockReturns struct {
		err error
	}

	type fields struct {
		aApp application.AuthenticationService
	}

	tests := []struct {
		name   string
		fields fields
		want
		mockArgs
		mockReturns
	}{
		{
			name: "When appropriate process completed, returns cookie which sessionID is empty and maxAge is 0",
			fields: fields{
				aApp: mock_application.NewMockAuthenticationService(ctrl),
			},
			want: want{
				statusCode: 200,
				cookie:     "SESSION_ID=; Path=/; HttpOnly",
				errBody: errBody{
					errCode: "",
					filed:   "",
				},
			},
			mockArgs: mockArgs{
				ctx:       context.Background(),
				sessionID: model.SessionValidIDForTest,
			},
			mockReturns: mockReturns{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aa, ok := tt.fields.aApp.(*mock_application.MockAuthenticationService)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}

			aa.EXPECT().Logout(tt.mockArgs.ctx, tt.mockArgs.sessionID).Return(tt.mockReturns.err)

			ac := NewAuthenticationController(tt.fields.aApp)
			r := gin.New()

			r.POST("/logout", ac.Logout)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/logout", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.AddCookie(&http.Cookie{
				Name:   model.SessionIDAtCookie,
				Value:  model.SessionValidIDForTest,
				MaxAge: 100,
			})

			r.ServeHTTP(rec, req)

			if rec.Code != tt.want.statusCode {
				t.Errorf("status code = %v, want %v", rec.Code, tt.want.statusCode)
				return
			}

			gotCookieVal := rec.Header().Get("Set-Cookie")
			if gotCookieVal != tt.want.cookie {
				t.Errorf("cookie = %v, want %v", gotCookieVal, tt.want.cookie)
				return
			}

			sBody := rec.Body.String()
			if !strings.Contains(sBody, string(tt.want.errBody.errCode)) {
				t.Errorf("body = %#v, want %#v", sBody, tt.want.errBody.errCode)
				return
			}

			if !strings.Contains(sBody, tt.want.errBody.filed) {
				t.Errorf("body = %#v, want %#v", sBody, tt.want.errBody.filed)
				return
			}
		})
	}
}
