package controller

import (
	"bytes"
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

func Test_commentController_ListComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		cApp application.CommentService
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
				cApp: mock_application.NewMockCommentService(ctrl),
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
				cApp: mock_application.NewMockCommentService(ctrl),
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
				cApp: mock_application.NewMockCommentService(ctrl),
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
				cApp: mock_application.NewMockCommentService(ctrl),
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
					DomainModelName: model.DomainModelNameComment,
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
			at, ok := tt.fields.cApp.(*mock_application.MockCommentService)
			if !ok {
				t.Fatal("failed to assert MockCommentService")
			}

			at.EXPECT().ListComments(context.Background(), tt.args.threadID, tt.args.limit, tt.args.cursor).Return(tt.mockReturns.list, tt.mockReturns.err)

			tc := NewCommentController(tt.fields.cApp)
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

func Test_commentController_GetComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		cApp application.CommentService
	}
	type args struct {
		threadID uint32
		id       uint32
	}

	type parameter struct {
		threadID string
		id       string
	}

	type errBody struct {
		errCode ErrCode
	}

	type want struct {
		statusCode int
		body       *model.Comment
		errBody
	}

	type mockReturns struct {
		comment *model.Comment
		err     error
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
			name: "When appropriate id is given and data exists, returns comment and status code 200",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				threadID: model.ThreadValidIDForTest,
				id:       model.CommentValidIDForTest,
			},
			parameter: parameter{
				threadID: "1",
				id:       "1",
			},
			mockReturns: mockReturns{
				comment: &model.Comment{
					ID:       model.CommentValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.Comment{
					ID:       model.CommentValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
				errBody: errBody{},
			},
		},
		{
			name: "When inappropriate id is given returns nil and status code 400",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			parameter: parameter{
				threadID: "1",
				id:       "a",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				errBody: errBody{
					errCode: InvalidParameterValueFailure,
				},
			},
		},
		{
			name: "When some error occurs, GetComment returns nil and status code 404",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				id: model.CommentInValidIDForTest,
			},
			parameter: parameter{
				threadID: "1",
				id:       "2",
			},
			mockReturns: mockReturns{
				comment: nil,
				err: &model.NoSuchDataError{
					PropertyName:  model.IDProperty,
					PropertyValue: "",
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
			at, ok := tt.fields.cApp.(*mock_application.MockCommentService)
			if !ok {
				t.Fatal("failed to assert MockCommentService")
			}

			if tt.want.statusCode != http.StatusBadRequest {
				at.EXPECT().GetComment(context.Background(), tt.args.id).Return(tt.mockReturns.comment, tt.mockReturns.err)
			}

			tc := NewCommentController(tt.fields.cApp)
			r := gin.New()

			r.GET("/threads/:threadId/comments/:id", tc.GetComment)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/threads/%s/comments/%s", tt.parameter.threadID, tt.parameter.id), nil)
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
				comment := &model.Comment{}
				if err := json.Unmarshal(bBody, comment); err != nil {
					t.Fatal(err)
					return
				}

				if !reflect.DeepEqual(comment, tt.want.body) {
					t.Errorf("body = %#v, want %#v", comment, tt.want.body)
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

func Test_commentController_CreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		cApp application.CommentService
	}
	type args struct {
		comment *CommentDTO
	}

	type errBody struct {
		errCode ErrCode
	}

	type want struct {
		statusCode int
		body       *model.Comment
		errBody
	}

	type mockArg struct {
		comment *model.Comment
	}

	type mockReturns struct {
		comment *model.Comment
		err     error
	}

	dto := &CommentDTO{
		ID:       model.CommentValidIDForTest,
		ThreadID: model.ThreadValidIDForTest,
		UserDTO: &UserDTO{
			ID:   model.UserValidIDForTest,
			Name: model.UserNameForTest,
		},
		Content: model.CommentContentForTest,
	}

	tests := []struct {
		name   string
		fields fields
		args
		mockArg
		mockReturns
		want
	}{
		{
			name: "When appropriate id is given and data exists, returns comment and status code 200",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				comment: dto,
			},
			mockArg: mockArg{
				comment: TranslateFromCommentDTOToComment(dto),
			},
			mockReturns: mockReturns{
				comment: &model.Comment{
					ID:       model.CommentValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.Comment{
					ID:       model.CommentValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
				errBody: errBody{},
			},
		},
		{
			name: "When validation error occurs, CreateComment returns nil and status code 400",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				comment: &CommentDTO{
					UserDTO: &UserDTO{},
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				errBody: errBody{
					errCode: InvalidParametersValueFailure,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			at, ok := tt.fields.cApp.(*mock_application.MockCommentService)
			if !ok {
				t.Fatal("failed to assert MockCommentService")
			}

			if tt.want.statusCode != http.StatusBadRequest {
				at.EXPECT().CreateComment(context.Background(), tt.mockArg.comment).Return(tt.mockReturns.comment, tt.mockReturns.err)
			}

			tc := NewCommentController(tt.fields.cApp)
			r := gin.New()

			r.POST("/threads/:threadId/comments", tc.CreateComment)

			rec := httptest.NewRecorder()

			b, err := json.Marshal(tt.args.comment)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/threads/1/comments", bytes.NewBuffer(b))
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
				comment := &model.Comment{}
				if err := json.Unmarshal(bBody, comment); err != nil {
					t.Fatal(err)
					return
				}

				if !reflect.DeepEqual(comment, tt.want.body) {
					t.Errorf("body = %#v, want %#v", comment, tt.want.body)
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

func Test_commentController_UpdateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		cApp application.CommentService
	}

	type parameter struct {
		id string
	}

	type args struct {
		id      uint32
		comment *CommentDTO
	}

	type errBody struct {
		errCode ErrCode
	}

	type want struct {
		statusCode int
		body       *model.Comment
		errBody
	}

	type mockArg struct {
		comment *model.Comment
	}

	type mockReturns struct {
		comment *model.Comment
		err     error
	}

	dto := &CommentDTO{
		ID:       model.CommentValidIDForTest,
		ThreadID: model.ThreadValidIDForTest,
		UserDTO: &UserDTO{
			ID:   model.UserValidIDForTest,
			Name: model.UserNameForTest,
		},
		Content: model.CommentContentForTest,
	}

	tests := []struct {
		name   string
		fields fields
		parameter
		args
		mockArg
		mockReturns
		want
	}{
		{
			name: "When appropriate id is given and data exists, returns comment and status code 200",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			parameter: parameter{
				id: "1",
			},
			args: args{
				id:      model.CommentValidIDForTest,
				comment: dto,
			},
			mockArg: mockArg{
				comment: TranslateFromCommentDTOToComment(dto),
			},
			mockReturns: mockReturns{
				comment: &model.Comment{
					ID:       model.CommentValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				body: &model.Comment{
					ID:       model.CommentValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
				errBody: errBody{},
			},
		},
		{
			name: "When validation error occurs, CreateComment returns nil and status code 400",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			parameter: parameter{
				id: "1",
			},
			args: args{
				id: model.CommentValidIDForTest,
				comment: &CommentDTO{
					UserDTO: &UserDTO{},
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				errBody: errBody{
					errCode: InvalidParametersValueFailure,
				},
			},
		},
		{
			name: "When inappropriate id is given returns nil and status code 400",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				id: model.CommentValidIDForTest,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			at, ok := tt.fields.cApp.(*mock_application.MockCommentService)
			if !ok {
				t.Fatal("failed to assert MockCommentService")
			}

			if tt.want.statusCode != http.StatusBadRequest {
				at.EXPECT().UpdateComment(context.Background(), tt.args.id, tt.mockArg.comment).Return(tt.mockReturns.comment, tt.mockReturns.err)
			}

			tc := NewCommentController(tt.fields.cApp)
			r := gin.New()

			r.PUT("/threads/:threadId/comments/:id", tc.UpdateComment)

			rec := httptest.NewRecorder()

			b, err := json.Marshal(tt.args.comment)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/threads/1/comments/%s", tt.parameter.id), bytes.NewBuffer(b))
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
				comment := &model.Comment{}
				if err := json.Unmarshal(bBody, comment); err != nil {
					t.Fatal(err)
					return
				}

				if !reflect.DeepEqual(comment, tt.want.body) {
					t.Errorf("body = %#v, want %#v", comment, tt.want.body)
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

func Test_commentController_DeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		cApp application.CommentService
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
		errBody
	}

	type mockReturns struct {
		err error
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
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				id: model.CommentValidIDForTest,
			},
			parameter: parameter{
				id: "1",
			},
			mockReturns: mockReturns{
				err: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				errBody:    errBody{},
			},
		},
		{
			name: "When inappropriate id is given returns nil and status code 400",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				id: model.CommentValidIDForTest,
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
			name: "When some error occurs, GetComment returns nil and status code 404",
			fields: fields{
				cApp: mock_application.NewMockCommentService(ctrl),
			},
			args: args{
				id: model.CommentValidIDForTest,
			},
			parameter: parameter{
				id: "1",
			},
			mockReturns: mockReturns{
				err: &model.NoSuchDataError{
					PropertyName:  model.IDProperty,
					PropertyValue: "",
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
			at, ok := tt.fields.cApp.(*mock_application.MockCommentService)
			if !ok {
				t.Fatal("failed to assert MockCommentService")
			}

			if tt.want.statusCode != http.StatusBadRequest {
				at.EXPECT().DeleteComment(context.Background(), tt.args.id).Return(tt.mockReturns.err)
			}

			tc := NewCommentController(tt.fields.cApp)
			r := gin.New()

			r.DELETE("/threads/:threadId/comments/:id", tc.DeleteComment)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/threads/1/comments/%s", tt.parameter.id), nil)
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
				thread := &model.Comment{}
				if err := json.Unmarshal(bBody, thread); err != nil {
					t.Fatal(err)
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
