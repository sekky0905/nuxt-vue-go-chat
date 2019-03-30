package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/pkg/errors"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

func TestNewCommentRepository(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want repository.CommentRepository
	}{
		{
			name: "When given appropriate args, returns commentRepository",
			args: args{
				ctx: context.Background(),
			},
			want: &commentRepository{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommentRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommentRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commentRepository_ErrorMsg(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		method model.RepositoryMethod
		err    error
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr *model.RepositoryError
	}{
		{
			name: "When given appropriate args, returns appropriate error",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				method: model.RepositoryMethodInsert,
				err:    errors.New(model.ErrorMessageForTest),
			},
			wantErr: &model.RepositoryError{
				BaseErr:                     errors.New(model.ErrorMessageForTest),
				RepositoryMethod:            model.RepositoryMethodInsert,
				DomainModelNameForDeveloper: model.DomainModelNameCommentForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameCommentForUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &commentRepository{}
			if err := repo.ErrorMsg(tt.args.method, tt.args.err); errors.Cause(err).Error() != tt.wantErr.Error() {
				t.Errorf("commentRepository{ErrorMsg() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}
