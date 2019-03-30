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

func Test_threadService_ListThreads(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		m        repository.DBManager
		repo     repository.ThreadRepository
		txCloser CloseTransaction
	}
	type args struct {
		ctx    context.Context
		limit  int
		cursor uint32
	}

	type mockReturns struct {
		list *model.ThreadList
		err  error
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		mockReturns
		want    *model.ThreadList
		wantErr bool
	}{
		{
			name: "When appropriate args given, ListThreads returns ThreadList and nil",
			fields: fields{
				m:    mock_repository.NewMockDBManager(ctrl),
				repo: mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx:    context.Background(),
				limit:  20,
				cursor: 21,
			},
			mockReturns: mockReturns{
				list: &model.ThreadList{
					Threads: testutil.GenerateThreadHelper(21, 40),
					HasNext: true,
					Cursor:  41,
				},
				err: nil,
			},
			want: &model.ThreadList{
				Threads: testutil.GenerateThreadHelper(21, 40),
				HasNext: true,
				Cursor:  41,
			},
			wantErr: false,
		},
		{
			name: "When some error occurs at repository layer, ListThreads returns nil and error",
			fields: fields{
				m:    mock_repository.NewMockDBManager(ctrl),
				repo: mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx:    context.Background(),
				limit:  20,
				cursor: 21,
			},
			mockReturns: mockReturns{
				list: nil,
				err:  errors.New(model.ErrorMessageForTest),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, ok := tt.fields.repo.(*mock_repository.MockThreadRepository)
			if !ok {
				t.Fatal("failed to assert MockThreadRepository")
			}
			tr.EXPECT().ListThreads(tt.args.ctx, tt.fields.m, tt.args.cursor, tt.args.limit).Return(tt.mockReturns.list, tt.mockReturns.err)

			a := &threadService{
				m:        tt.fields.m,
				repo:     tt.fields.repo,
				txCloser: tt.fields.txCloser,
			}
			got, err := a.ListThreads(tt.args.ctx, tt.args.limit, tt.args.cursor)
			if (err != nil) != tt.wantErr {
				t.Errorf("threadService.ListThreads() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("threadService.ListThreads() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_threadService_GetThread(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		m        repository.DBManager
		repo     repository.ThreadRepository
		txCloser CloseTransaction
	}
	type args struct {
		ctx context.Context
		id  uint32
	}

	type mockReturns struct {
		thread *model.Thread
		err    error
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		mockReturns
		want    *model.Thread
		wantErr bool
	}{
		{
			name: "When appropriate args given, GetThread returns Thread and nil",
			fields: fields{
				m:    mock_repository.NewMockDBManager(ctrl),
				repo: mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockReturns: mockReturns{
				thread: &model.Thread{
					ID:    uint32(model.ThreadValidIDForTest),
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
			want: &model.Thread{
				ID:    uint32(model.ThreadValidIDForTest),
				Title: model.TitleForTest,
				User: &model.User{
					ID:   model.UserValidIDForTest,
					Name: model.UserNameForTest,
				},
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: false,
		},
		{
			name: "When some error occurs at repository layer, GetThread returns nil and error",
			fields: fields{
				m:    mock_repository.NewMockDBManager(ctrl),
				repo: mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				id:  model.ThreadInValidIDForTest,
			},
			mockReturns: mockReturns{
				thread: nil,
				err:    errors.New(model.ErrorMessageForTest),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, ok := tt.fields.repo.(*mock_repository.MockThreadRepository)
			if !ok {
				t.Fatal("failed to assert MockThreadRepository")
			}
			tr.EXPECT().GetThreadByID(tt.args.ctx, tt.fields.m, tt.args.id).Return(tt.mockReturns.thread, tt.mockReturns.err)

			a := &threadService{
				m:        tt.fields.m,
				repo:     tt.fields.repo,
				txCloser: tt.fields.txCloser,
			}
			got, err := a.GetThread(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("threadService.GetThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("threadService.GetThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_threadService_CreateThread(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		m        repository.DBManager
		service  service.ThreadService
		repo     repository.ThreadRepository
		txCloser CloseTransaction
	}
	type args struct {
		ctx   context.Context
		param *model.Thread
	}

	type mockArgsIsAlreadyExistTitle struct {
		ctx   context.Context
		title string
	}

	type mockReturnsIsAlreadyExistTitle struct {
		found bool
		err   error
	}

	type mockArgsInsertThread struct {
		ctx   context.Context
		tx    repository.DBManager
		param *model.Thread
	}

	type mockReturnsInsertThread struct {
		id  uint32
		err error
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		mockArgsIsAlreadyExistTitle
		mockReturnsIsAlreadyExistTitle
		mockArgsInsertThread
		mockReturnsInsertThread
		wantThread *model.Thread
		wantErr    bool
	}{
		{
			name: "When appropriate args given, CreateThread returns id and nil",
			fields: fields{
				m:       mock_repository.NewMockDBManager(ctrl),
				repo:    mock_repository.NewMockThreadRepository(ctrl),
				service: mock_service.NewMockThreadService(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				param: &model.Thread{
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockArgsIsAlreadyExistTitle: mockArgsIsAlreadyExistTitle{
				ctx:   context.Background(),
				title: model.TitleForTest,
			},
			mockReturnsIsAlreadyExistTitle: mockReturnsIsAlreadyExistTitle{
				found: false,
				err:   &model.NoSuchDataError{},
			},
			mockArgsInsertThread: mockArgsInsertThread{
				ctx: context.Background(),
				param: &model.Thread{
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockReturnsInsertThread: mockReturnsInsertThread{
				id:  model.ThreadValidIDForTest,
				err: nil,
			},
			wantThread: &model.Thread{
				ID:    model.ThreadValidIDForTest,
				Title: model.TitleForTest,
				User: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: false,
		},
		{
			name: "When given id has already existed, CreateThread returns nil and error",
			fields: fields{
				m:       mock_repository.NewMockDBManager(ctrl),
				repo:    mock_repository.NewMockThreadRepository(ctrl),
				service: mock_service.NewMockThreadService(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				param: &model.Thread{
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockArgsIsAlreadyExistTitle: mockArgsIsAlreadyExistTitle{
				ctx:   context.Background(),
				title: model.TitleForTest,
			},
			mockReturnsIsAlreadyExistTitle: mockReturnsIsAlreadyExistTitle{
				found: true,
				err:   nil,
			},
			mockArgsInsertThread: mockArgsInsertThread{
				ctx:   context.Background(),
				param: nil,
			},
			wantThread: nil,
			wantErr:    true,
		},
		{
			name: "When some error occurs at repository layer, CreateThread returns nil and error",
			fields: fields{
				m:       mock_repository.NewMockDBManager(ctrl),
				repo:    mock_repository.NewMockThreadRepository(ctrl),
				service: mock_service.NewMockThreadService(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				param: &model.Thread{
					ID:    uint32(model.ThreadValidIDForTest),
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockArgsIsAlreadyExistTitle: mockArgsIsAlreadyExistTitle{
				ctx:   context.Background(),
				title: model.TitleForTest,
			},
			mockReturnsIsAlreadyExistTitle: mockReturnsIsAlreadyExistTitle{
				found: false,
				err:   &model.NoSuchDataError{},
			},
			mockArgsInsertThread: mockArgsInsertThread{
				ctx: context.Background(),
				tx:  mock_repository.NewMockDBManager(ctrl),
				param: &model.Thread{
					ID:    uint32(model.ThreadValidIDForTest),
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockReturnsInsertThread: mockReturnsInsertThread{
				id:  model.ThreadValidIDForTest,
				err: errors.New(model.ErrorMessageForTest),
			},
			wantThread: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, ok := tt.fields.m.(*mock_repository.MockDBManager)
			if !ok {
				t.Fatal("failed to assert MockDBManager")
			}
			m.EXPECT().Begin().Return(mock_repository.NewMockTxManager(ctrl), nil)

			ts, ok := tt.fields.service.(*mock_service.MockThreadService)
			if !ok {
				t.Fatal("failed to assert MockThreadService")
			}

			ts.EXPECT().IsAlreadyExistTitle(tt.mockArgsIsAlreadyExistTitle.ctx, tt.mockArgsIsAlreadyExistTitle.title).Return(tt.mockReturnsIsAlreadyExistTitle.found, tt.mockReturnsIsAlreadyExistTitle.err)

			if tt.mockArgsInsertThread.param != nil {
				tr, ok := tt.fields.repo.(*mock_repository.MockThreadRepository)
				if !ok {
					t.Fatal("failed to assert MockThreadRepository")
				}

				txM := mock_repository.NewMockTxManager(ctrl)

				tr.EXPECT().InsertThread(tt.mockArgsInsertThread.ctx, txM, tt.args.param).Return(tt.mockReturnsInsertThread.id, tt.mockReturnsInsertThread.err)
			}

			a := &threadService{
				m:        tt.fields.m,
				repo:     tt.fields.repo,
				service:  tt.fields.service,
				txCloser: tt.fields.txCloser,
			}
			gotThread, err := a.CreateThread(tt.args.ctx, tt.args.param)
			if gotThread != nil {
				gotThread.CreatedAt = tt.wantThread.CreatedAt
				gotThread.UpdatedAt = tt.wantThread.UpdatedAt
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("threadService.CreateThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotThread, tt.wantThread) {
				t.Errorf("threadService.CreateThread() = %v, want %v", gotThread, tt.wantThread)
			}
		})
	}
}

func Test_threadService_UpdateThread(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		m        repository.DBManager
		service  service.ThreadService
		repo     repository.ThreadRepository
		txCloser CloseTransaction
	}
	type args struct {
		ctx   context.Context
		id    uint32
		param *model.Thread
	}

	type mockArgsIsAlreadyExistID struct {
		ctx context.Context
		id  uint32
	}

	type mockReturnsIsAlreadyExistID struct {
		found bool
		err   error
	}

	type mockArgsUpdateThread struct {
		ctx   context.Context
		param *model.Thread
	}

	type mockReturnsUpdateThread struct {
		err error
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		mockArgsIsAlreadyExistID
		mockReturnsIsAlreadyExistID
		mockArgsUpdateThread
		mockReturnsUpdateThread
		wantThread *model.Thread
		wantErr    bool
	}{
		{
			name: "When appropriate args given, UpdateThread returns Thread and err",
			fields: fields{
				m:       mock_repository.NewMockDBManager(ctrl),
				service: mock_service.NewMockThreadService(ctrl),
				repo:    mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
				param: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockArgsIsAlreadyExistID: mockArgsIsAlreadyExistID{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
			},
			mockReturnsIsAlreadyExistID: mockReturnsIsAlreadyExistID{
				found: true,
				err:   nil,
			},
			mockArgsUpdateThread: mockArgsUpdateThread{
				ctx: context.Background(),
				param: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockReturnsUpdateThread: mockReturnsUpdateThread{
				err: nil,
			},
			wantThread: &model.Thread{
				ID:    model.ThreadValidIDForTest,
				Title: model.TitleForTest,
				User: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: false,
		},
		{
			name: "When given id has not existed, UpdateThread returns nil and error",
			fields: fields{
				m:       mock_repository.NewMockDBManager(ctrl),
				service: mock_service.NewMockThreadService(ctrl),
				repo:    mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
				param: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockArgsIsAlreadyExistID: mockArgsIsAlreadyExistID{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
			},
			mockReturnsIsAlreadyExistID: mockReturnsIsAlreadyExistID{
				found: false,
				err:   nil,
			},
			mockArgsUpdateThread: mockArgsUpdateThread{
				param: nil,
			},
			wantThread: nil,
			wantErr:    true,
		},
		{
			name: "When some error occurs at repository layer, UpdateThread returns nil and error",
			fields: fields{
				m:       mock_repository.NewMockDBManager(ctrl),
				service: mock_service.NewMockThreadService(ctrl),
				repo:    mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
				param: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockArgsIsAlreadyExistID: mockArgsIsAlreadyExistID{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
			},
			mockReturnsIsAlreadyExistID: mockReturnsIsAlreadyExistID{
				found: true,
				err:   nil,
			},
			mockArgsUpdateThread: mockArgsUpdateThread{
				ctx: context.Background(),
				param: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockReturnsUpdateThread: mockReturnsUpdateThread{
				err: errors.New(model.ErrorMessageForTest),
			},
			wantThread: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, ok := tt.fields.m.(*mock_repository.MockDBManager)
			if !ok {
				t.Fatal("failed to assert MockDBManager")
			}
			m.EXPECT().Begin().Return(mock_repository.NewMockTxManager(ctrl), nil)

			ts, ok := tt.fields.service.(*mock_service.MockThreadService)
			if !ok {
				t.Fatal("failed to assert MockThreadService")
			}

			ts.EXPECT().IsAlreadyExistID(tt.mockArgsIsAlreadyExistID.ctx, tt.mockArgsIsAlreadyExistID.id).Return(tt.mockReturnsIsAlreadyExistID.found, tt.mockReturnsIsAlreadyExistID.err)

			if tt.mockArgsUpdateThread.param != nil {

				tr, ok := tt.fields.repo.(*mock_repository.MockThreadRepository)
				if !ok {
					t.Fatal("failed to assert MockThreadRepository")
				}

				txM := mock_repository.NewMockTxManager(ctrl)

				tr.EXPECT().UpdateThread(tt.mockArgsUpdateThread.ctx, txM, tt.args.id, tt.args.param).Return(tt.mockReturnsUpdateThread.err)

			}

			a := &threadService{
				m:        tt.fields.m,
				service:  tt.fields.service,
				repo:     tt.fields.repo,
				txCloser: tt.fields.txCloser,
			}

			gotThread, err := a.UpdateThread(tt.args.ctx, tt.args.id, tt.args.param)
			if gotThread != nil {
				gotThread.UpdatedAt = tt.wantThread.UpdatedAt
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("threadService.UpdateThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotThread, tt.wantThread) {
				t.Errorf("threadService.UpdateThread() = %+v, want %+v", gotThread, tt.wantThread)
			}
		})
	}
}

func Test_threadService_DeleteThread(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		m        repository.DBManager
		service  service.ThreadService
		repo     repository.ThreadRepository
		txCloser CloseTransaction
	}
	type args struct {
		ctx   context.Context
		id    uint32
		param *model.Thread
	}

	type mockArgsIsAlreadyExistID struct {
		ctx context.Context
		id  uint32
	}

	type mockReturnsIsAlreadyExistID struct {
		found bool
		err   error
	}

	type mockReturnsDeleteThread struct {
		err error
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		mockArgsIsAlreadyExistID
		mockReturnsIsAlreadyExistID
		mockReturnsDeleteThread
		wantThread *model.Thread
		wantErr    bool
	}{
		{
			name: "When appropriate args given, DeleteThread returns Thread and err",
			fields: fields{
				m:       mock_repository.NewMockDBManager(ctrl),
				service: mock_service.NewMockThreadService(ctrl),
				repo:    mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
				param: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockArgsIsAlreadyExistID: mockArgsIsAlreadyExistID{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
			},
			mockReturnsIsAlreadyExistID: mockReturnsIsAlreadyExistID{
				found: true,
				err:   nil,
			},
			mockReturnsDeleteThread: mockReturnsDeleteThread{
				err: nil,
			},
			wantThread: &model.Thread{
				ID:    model.ThreadValidIDForTest,
				Title: model.TitleForTest,
				User: &model.User{
					ID:        model.UserValidIDForTest,
					Name:      model.UserNameForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: false,
		},
		{
			name: "When given id has not existed, DeleteThread returns nil and error",
			fields: fields{
				m:       mock_repository.NewMockDBManager(ctrl),
				service: mock_service.NewMockThreadService(ctrl),
				repo:    mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
				param: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockArgsIsAlreadyExistID: mockArgsIsAlreadyExistID{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
			},
			mockReturnsIsAlreadyExistID: mockReturnsIsAlreadyExistID{
				found: false,
				err:   nil,
			},
			wantThread: nil,
			wantErr:    true,
		},
		{
			name: "When some error occurs at repository layer, DeleteThread returns nil and error",
			fields: fields{
				m:       mock_repository.NewMockDBManager(ctrl),
				service: mock_service.NewMockThreadService(ctrl),
				repo:    mock_repository.NewMockThreadRepository(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
				param: &model.Thread{
					ID:    model.ThreadValidIDForTest,
					Title: model.TitleForTest,
					User: &model.User{
						ID:        model.UserValidIDForTest,
						Name:      model.UserNameForTest,
						CreatedAt: testutil.TimeNow(),
						UpdatedAt: testutil.TimeNow(),
					},
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockArgsIsAlreadyExistID: mockArgsIsAlreadyExistID{
				ctx: context.Background(),
				id:  model.ThreadValidIDForTest,
			},
			mockReturnsIsAlreadyExistID: mockReturnsIsAlreadyExistID{
				found: true,
				err:   nil,
			},
			mockReturnsDeleteThread: mockReturnsDeleteThread{
				err: errors.New(model.ErrorMessageForTest),
			},
			wantThread: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, ok := tt.fields.m.(*mock_repository.MockDBManager)
			if !ok {
				t.Fatal("failed to assert MockDBManager")
			}
			m.EXPECT().Begin().Return(mock_repository.NewMockTxManager(ctrl), nil)

			ts, ok := tt.fields.service.(*mock_service.MockThreadService)
			if !ok {
				t.Fatal("failed to assert MockThreadService")
			}

			ts.EXPECT().IsAlreadyExistID(tt.mockArgsIsAlreadyExistID.ctx, tt.mockArgsIsAlreadyExistID.id).Return(tt.mockReturnsIsAlreadyExistID.found, tt.mockReturnsIsAlreadyExistID.err)

			if tt.mockReturnsIsAlreadyExistID.found {
				tr, ok := tt.fields.repo.(*mock_repository.MockThreadRepository)
				if !ok {
					t.Fatal("failed to assert MockThreadRepository")
				}

				txM := mock_repository.NewMockTxManager(ctrl)

				tr.EXPECT().DeleteThread(tt.args.ctx, txM, tt.args.id).Return(tt.mockReturnsDeleteThread.err)
			}

			a := &threadService{
				m:        tt.fields.m,
				service:  tt.fields.service,
				repo:     tt.fields.repo,
				txCloser: tt.fields.txCloser,
			}

			err := a.DeleteThread(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("threadService.DeleteThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
