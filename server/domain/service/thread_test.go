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
	"github.com/sekky0905/nuxt-vue-go-chat/server/testutil"
)

func Test_threadService_IsAlreadyExistID(t *testing.T) {
	// for gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock_repository.NewMockThreadRepository(ctrl)

	type fields struct {
		repo repository.ThreadRepository
	}
	type args struct {
		ctx context.Context
		m   query.SQLManager
		id  uint32
	}

	type returnArgs struct {
		thread *model.Thread
		err    error
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
			name: "When specified thread already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx: context.Background(),
				m:   db.NewDBManager(),
				id:  model.ThreadValidIDForTest,
			},
			returnArgs: returnArgs{
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: nil,
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "When specified thread doesn't already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx: context.Background(),
				m:   db.NewDBManager(),
				id:  model.ThreadInValidIDForTest,
			},
			returnArgs: returnArgs{
				thread: nil,
				err:    nil,
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
				id:  model.ThreadInValidIDForTest,
			},
			returnArgs: returnArgs{
				thread: nil,
				err:    errors.New(model.ErrorMessageForTest),
			},
			want:    false,
			wantErr: errors.New(model.ErrorMessageForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &threadService{
				repo: mock,
			}

			mock.EXPECT().GetThreadByID(tt.args.ctx, tt.args.m, tt.args.id).Return(tt.returnArgs.thread, tt.returnArgs.err)

			got, err := s.IsAlreadyExistID(tt.args.ctx, tt.args.m, tt.args.id)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("threadService.IsAlreadyExistID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("threadService.IsAlreadyExistID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_threadService_IsAlreadyExistTitle(t *testing.T) {
	// for gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock_repository.NewMockThreadRepository(ctrl)

	type fields struct {
		repo repository.ThreadRepository
	}
	type args struct {
		ctx   context.Context
		m     query.SQLManager
		title string
	}

	type returnArgs struct {
		thread *model.Thread
		err    error
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
			name: "When specified thread already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx:   context.Background(),
				m:     db.NewDBManager(),
				title: model.TitleForTest,
			},
			returnArgs: returnArgs{
				thread: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:   model.UserValidIDForTest,
						Name: model.UserNameForTest,
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: nil,
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "When specified thread doesn't already exists, return true and nil.",
			fields: fields{
				repo: mock,
			},
			args: args{
				ctx:   context.Background(),
				m:     db.NewDBManager(),
				title: model.TitleForTest,
			},
			returnArgs: returnArgs{
				thread: nil,
				err:    nil,
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
				ctx:   context.Background(),
				m:     db.NewDBManager(),
				title: model.TitleForTest,
			},
			returnArgs: returnArgs{
				thread: nil,
				err:    errors.New(model.ErrorMessageForTest),
			},
			want:    false,
			wantErr: errors.New(model.ErrorMessageForTest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &threadService{
				repo: mock,
			}

			mock.EXPECT().GetThreadByTitle(tt.args.ctx, tt.args.m, tt.args.title).Return(tt.returnArgs.thread, tt.returnArgs.err)

			got, err := s.IsAlreadyExistTitle(tt.args.ctx, tt.args.m, tt.args.title)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("threadService.IsAlreadyExistTitle() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("threadService.IsAlreadyExistTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
