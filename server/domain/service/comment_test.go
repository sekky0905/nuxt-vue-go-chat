package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	mock_repository "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository/mock"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"
)

func Test_commentService_IsAlreadyExistID(t *testing.T) {
	// for gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock_repository.NewMockCommentRepository(ctrl)

	type fields struct {
		repo repository.CommentRepository
	}
	type args struct {
		ctx context.Context
		m   query.SQLManager
		id  uint32
	}

	type returnArgs struct {
		comment *model.Comment
		err     error
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
			name: "When specified comment already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx: context.Background(),
				m:   db.NewDBManager(),
				id:  model.CommentValidIDForTest,
			},
			returnArgs: returnArgs{
				comment: &model.Comment{
					ID:       model.CommentInValidIDForTest,
					ThreadID: model.ThreadValidIDForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					Content: model.CommentContentForTest,
				},
				err: nil,
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "When specified comment doesn't already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx: context.Background(),
				m:   db.NewDBManager(),
				id:  model.CommentInValidIDForTest,
			},
			returnArgs: returnArgs{
				comment: nil,
				err:     nil,
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
				id:  model.CommentInValidIDForTest,
			},
			returnArgs: returnArgs{
				comment: nil,
				err:     errors.New(model.ErrorMessageForTest),
			},
			want:    false,
			wantErr: errors.New(model.ErrorMessageForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &commentService{
				repo: mock,
			}

			mock.EXPECT().GetCommentByID(tt.args.ctx, tt.args.m, tt.args.id).Return(tt.returnArgs.comment, tt.returnArgs.err)

			got, err := s.IsAlreadyExistID(tt.args.ctx, tt.args.m, tt.args.id)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("commentService.IsAlreadyExistID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("commentService.IsAlreadyExistID() = %v, want %v", got, tt.want)
			}
		})
	}
}
