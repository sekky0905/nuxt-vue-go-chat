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

	"github.com/golang/mock/gomock"
	mock_application "github.com/sekky0905/nuxt-vue-go-chat/server/application/mock"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"

	"github.com/gin-gonic/gin"
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
)

func Test_commentController_ListComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		tApp application.CommentService
	}
	type args struct {
		threadID uint32
		limit    int
		cursor   uint32
	}

	type parameter struct {
		threadID string
		limit    string
		cursor   string
	}

	type errBody struct {
		errCode ErrCode
	}

	type want struct {
		statusCode int
		body       *model.CommentList
		errBody
	}

	type mockReturns struct {
		list *model.CommentList
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
			name: "When appropriate limit and commit are given and data exists, returns commentList and status code 200",
			fields: fields{
				tApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				threadID: model.ThreadValidIDForTest,
				limit:    20,
				cursor:   21,
			},
			parameter: parameter{
				threadID: string(model.ThreadValidIDForTest),
				limit:    "20",
				cursor:   "21",
			},
			mockReturns: mockReturns{
				list: &model.CommentList{
					Comments: testutil.GenerateCommentHelper(21, 40),
					HasNext:  true,
					Cursor:   41,
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.CommentList{
					Comments: testutil.GenerateCommentHelper(21, 40),
					HasNext:  true,
					Cursor:   41,
				},
				errBody: errBody{},
			},
		},
		{
			name: "When inappropriate limit is given and data exists, returns commentList which has 20 data and status code 200",
			fields: fields{
				tApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				threadID: model.ThreadValidIDForTest,
				limit:    20,
				cursor:   21,
			},
			parameter: parameter{
				threadID: string(model.ThreadValidIDForTest),
				limit:    "test",
				cursor:   "21",
			},
			mockReturns: mockReturns{
				list: &model.CommentList{
					Comments: testutil.GenerateCommentHelper(21, 40),
					HasNext:  true,
					Cursor:   41,
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.CommentList{
					Comments: testutil.GenerateCommentHelper(21, 40),
					HasNext:  true,
					Cursor:   41,
				},
				errBody: errBody{},
			},
		},
		{
			name: "When inappropriate cursor is given and data exists, returns commentList which has 1~21 data and status code 200",
			fields: fields{
				tApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				threadID: model.ThreadValidIDForTest,
				limit:    20,
				cursor:   1,
			},
			parameter: parameter{
				threadID: string(model.ThreadValidIDForTest),
				limit:    "20",
				cursor:   "test",
			},
			mockReturns: mockReturns{
				list: &model.CommentList{
					Comments: testutil.GenerateCommentHelper(1, 20),
					HasNext:  true,
					Cursor:   21,
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.CommentList{
					Comments: testutil.GenerateCommentHelper(1, 20),
					HasNext:  true,
					Cursor:   21,
				},
				errBody: errBody{},
			},
		},
		{
			name: "When inappropriate cursor is given and data exists, returns error status code 404",
			fields: fields{
				tApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				threadID: model.ThreadValidIDForTest,
				limit:    20,
				cursor:   1,
			},
			parameter: parameter{
				threadID: string(model.ThreadValidIDForTest),
				limit:    "20",
				cursor:   "test",
			},
			mockReturns: mockReturns{
				list: nil,
				err: &model.NoSuchDataError{
					DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
					DomainModelNameForUser:      model.DomainModelNameCommentForUser,
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
			at, ok := tt.fields.tApp.(*mock_application.MockCommentService)
			if !ok {
				t.Fatal("failed to assert MockCommentService")
			}

			at.EXPECT().ListComments(context.Background(), tt.args.threadID, tt.args.limit, tt.args.cursor).Return(tt.mockReturns.list, tt.mockReturns.err)

			tc := NewCommentController(tt.fields.tApp)
			r := gin.New()

			r.GET("/threads/:threadId/comments", tc.ListComments)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/threads/%d/comments?limit=%s&cursor=%s", tt.args.threadID, tt.parameter.limit, tt.parameter.cursor), nil)
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
				commentList := &model.CommentList{}
				if err := json.Unmarshal(bBody, commentList); err != nil {
					t.Fatal(err)
					return
				}

				if !reflect.DeepEqual(commentList.Comments, tt.want.body.Comments) {
					t.Errorf("body = %+v, want %+v", commentList.Comments, tt.want.body.Comments)
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
