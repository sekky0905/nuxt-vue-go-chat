package controller

import (
	"context"
	"encoding/json"
	"fmt"
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

func Test_threadController_ListThreads(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		tApp application.ThreadService
	}
	type args struct {
		limit  int
		cursor uint32
	}

	type parameter struct {
		limit  string
		cursor string
	}

	type errBody struct {
		errCode ErrCode
	}

	type want struct {
		statusCode int
		body       *model.ThreadList
		errBody
	}

	type mockReturns struct {
		list *model.ThreadList
		err  error
	}

	tests := []struct {
		name   string
		fields fields
		args
		parameter
		mockReturns
		want
	}{
		{
			name: "When appropriate limit and commit are given and data exists, returns threadList and status code 200",
			fields: fields{
				tApp: mock_application.NewMockThreadService(ctrl),
			},
			args: args{
				limit:  20,
				cursor: 21,
			},
			parameter: parameter{
				limit:  "20",
				cursor: "21",
			},
			mockReturns: mockReturns{
				list: &model.ThreadList{
					Threads: testutil.GenerateThreadHelper(21, 40),
					HasNext: true,
					Cursor:  41,
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.ThreadList{
					Threads: testutil.GenerateThreadHelper(21, 40),
					HasNext: true,
					Cursor:  41,
				},
				errBody: errBody{},
			},
		},
		{
			name: "When inappropriate limit is given and data exists, returns threadList which has 20 data and status code 200",
			fields: fields{
				tApp: mock_application.NewMockThreadService(ctrl),
			},
			args: args{
				limit:  20,
				cursor: 21,
			},
			parameter: parameter{
				limit:  "test",
				cursor: "21",
			},
			mockReturns: mockReturns{
				list: &model.ThreadList{
					Threads: testutil.GenerateThreadHelper(21, 40),
					HasNext: true,
					Cursor:  41,
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.ThreadList{
					Threads: testutil.GenerateThreadHelper(21, 40),
					HasNext: true,
					Cursor:  41,
				},
				errBody: errBody{},
			},
		},
		{
			name: "When inappropriate cursor is given and data exists, returns threadList which has 1~21 data and status code 200",
			fields: fields{
				tApp: mock_application.NewMockThreadService(ctrl),
			},
			args: args{
				limit:  20,
				cursor: 1,
			},
			parameter: parameter{
				limit:  "20",
				cursor: "test",
			},
			mockReturns: mockReturns{
				list: &model.ThreadList{
					Threads: testutil.GenerateThreadHelper(1, 20),
					HasNext: true,
					Cursor:  21,
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.ThreadList{
					Threads: testutil.GenerateThreadHelper(1, 20),
					HasNext: true,
					Cursor:  21,
				},
				errBody: errBody{},
			},
		},
		{
			name: "When inappropriate cursor is given and data exists, returns error status code 404",
			fields: fields{
				tApp: mock_application.NewMockThreadService(ctrl),
			},
			args: args{
				limit:  20,
				cursor: 1,
			},
			parameter: parameter{
				limit:  "20",
				cursor: "test",
			},
			mockReturns: mockReturns{
				list: nil,
				err: &model.NoSuchDataError{
					DomainModelNameForDeveloper: model.DomainModelNameThreadForDeveloper,
					DomainModelNameForUser:      model.DomainModelNameThreadForUser,
				},
			},
			want: want{
				statusCode: http.StatusNotFound,
				body:       nil,
				errBody: errBody{
					errCode: NoSuchDataFailure,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			at, ok := tt.fields.tApp.(*mock_application.MockThreadService)
			if !ok {
				t.Fatal("failed to assert MockThreadService")
			}

			at.EXPECT().ListThreads(context.Background(), tt.args.limit, tt.args.cursor).Return(tt.mockReturns.list, tt.mockReturns.err)

			tc := NewThreadController(tt.fields.tApp)
			r := gin.New()

			r.GET("/threads", tc.ListThreads)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/threads?limit=%s&cursor=%s", tt.parameter.limit, tt.parameter.cursor), nil)
			if err != nil {
				t.Fatal(err)
			}
			r.ServeHTTP(rec, req)

			if rec.Code != tt.want.statusCode {
				t.Errorf("status code = %v, want %v", rec.Code, tt.want.statusCode)
				return
			}

			if tt.want.errBody.errCode == "" {
				bBody := rec.Body.Bytes()
				threadList := &model.ThreadList{}
				if err := json.Unmarshal(bBody, threadList); err != nil {
					t.Fatal(err)
					return
				}

				if !reflect.DeepEqual(threadList.Threads, tt.want.body.Threads) {
					t.Errorf("body = %+v, want %+v", threadList.Threads, tt.want.body.Threads)
					return
				}
			} else {
				sBody := rec.Body.String()
				if !strings.Contains(sBody, string(tt.want.errBody.errCode)) {
					t.Errorf("body = %#v, want %#v", sBody, tt.want.errBody.errCode)
					return
				}
			}
		})
	}
}

func Test_threadController_GetThread(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		tApp application.ThreadService
	}
	type args struct {
		id uint32
	}

	type parameter struct {
		id string
	}

	type errBody struct {
		errCode ErrCode
	}

	type want struct {
		statusCode int
		body       *model.Thread
		errBody
	}

	type mockReturns struct {
		thread *model.Thread
		err    error
	}

	tests := []struct {
		name   string
		fields fields
		args
		parameter
		mockReturns
		want
	}{
		{
			name: "When appropriate id is given and data exists, returns thread and status code 200",
			fields: fields{
				tApp: mock_application.NewMockThreadService(ctrl),
			},
			args: args{
				id: model.ThreadValidIDForTest,
			},
			parameter: parameter{
				id: "1",
			},
			mockReturns: mockReturns{
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
				},
				errBody: errBody{},
			},
		},
		{
			name: "When inappropriate id is given returns nil and status code 400",
			fields: fields{
				tApp: mock_application.NewMockThreadService(ctrl),
			},
			args: args{
				id: model.ThreadValidIDForTest,
			},
			parameter: parameter{
				id: "a",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				errBody: errBody{
					errCode: InvalidParametersValueFailure,
				},
			},
		},
		{
			name: "When some error occurs, GetThread returns nil and status code 404",
			fields: fields{
				tApp: mock_application.NewMockThreadService(ctrl),
			},
			args: args{
				id: model.ThreadValidIDForTest,
			},
			parameter: parameter{
				id: "1",
			},
			mockReturns: mockReturns{
				thread: nil,
				err: &model.NoSuchDataError{
					PropertyNameForDeveloper: model.IDPropertyForDeveloper,
					PropertyValue:            "",
				},
			},
			want: want{
				statusCode: http.StatusNotFound,
				errBody: errBody{
					errCode: NoSuchDataFailure,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			at, ok := tt.fields.tApp.(*mock_application.MockThreadService)
			if !ok {
				t.Fatal("failed to assert MockThreadService")
			}

			if tt.want.statusCode != http.StatusBadRequest {
				at.EXPECT().GetThread(context.Background(), tt.args.id).Return(tt.mockReturns.thread, tt.mockReturns.err)
			}

			tc := NewThreadController(tt.fields.tApp)
			r := gin.New()

			r.GET("/threads/:id", tc.GetThread)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/threads/%s", tt.parameter.id), nil)
			if err != nil {
				t.Fatal(err)
			}
			r.ServeHTTP(rec, req)

			if rec.Code != tt.want.statusCode {
				t.Errorf("status code = %v, want %v", rec.Code, tt.want.statusCode)
				return
			}

			if tt.want.errBody.errCode == "" {
				bBody := rec.Body.Bytes()
				thread := &model.Thread{}
				if err := json.Unmarshal(bBody, thread); err != nil {
					t.Fatal(err)
					return
				}

				if !reflect.DeepEqual(thread, tt.want.body) {
					t.Errorf("body = %#v, want %#v", thread, tt.want.body)
					return
				}
			} else {
				sBody := rec.Body.String()
				if !strings.Contains(sBody, string(tt.want.errBody.errCode)) {
					t.Errorf("body = %#v, want %#v", sBody, tt.want.errBody.errCode)
					return
				}
			}
		})
	}
}
