package db

import (
	"context"
	"reflect"
	"testing"

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
