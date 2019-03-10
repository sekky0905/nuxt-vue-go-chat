package application

import (
	"context"
	"reflect"
	"testing"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
)

func Test_authenticationService_SignUp(t *testing.T) {
	type fields struct {
		m        repository.SQLManager
		uFactory service.UserRepoFactory
		sFactory service.SessionRepoFactory
		txCloser CloseTransaction
	}
	type args struct {
		ctx      context.Context
		name     string
		password string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser *model.User
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authenticationService{
				m:        tt.fields.m,
				uFactory: tt.fields.uFactory,
				sFactory: tt.fields.sFactory,
				txCloser: tt.fields.txCloser,
			}
			gotUser, err := a.SignUp(tt.args.ctx, tt.args.name, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("authenticationService.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("authenticationService.SignUp() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
