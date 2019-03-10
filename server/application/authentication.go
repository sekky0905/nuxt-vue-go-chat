package application

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
	"github.com/sekky0905/nuxt-vue-go-chat/server/util"
)

// AuthenticationService is the interface of AuthenticationService.
type AuthenticationService interface {
	SignUp(ctx context.Context, name, password string) (*model.User, error)
}

// authenticationService is the service of authentication.
type authenticationService struct {
	m        repository.SQLManager
	uFactory service.UserRepoFactory
	sFactory service.SessionRepoFactory
	txCloser CloseTransaction
}

// NewAuthenticationApplication generates and returns AuthenticationService.
func NewAuthenticationApplication(m repository.SQLManager, uFactory service.UserRepoFactory, sFactory service.SessionRepoFactory, txCloser CloseTransaction) AuthenticationService {
	return &authenticationService{
		m:        m,
		uFactory: uFactory,
		sFactory: sFactory,
		txCloser: txCloser,
	}
}

// SignUp sign up an user.
func (a *authenticationService) SignUp(ctx context.Context, name, password string) (user *model.User, err error) {
	tx, err := a.m.Begin()
	if err != nil {
		return nil, beginTxErrorMsg(err)
	}

	defer func() {
		if err := a.txCloser(tx, err); err != nil {
			err = errors.Wrap(err, "failed to close tx")
		}
	}()

	user, err = model.NewUser(name, password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new user")
	}

	uRepo := a.uFactory(ctx)
	uService := service.NewUserService(uRepo, a.m)

	// get User
	if existingUser, err := uRepo.GetUserByName(a.m, name); existingUser != nil {
		err = &model.AuthenticationErr{
			BaseErr: err,
		}
		return nil, errors.Wrapf(err, "already registered")
	} else if err != nil && errors.Cause(err) != err.(*model.NoSuchDataError) {
		return nil, errors.Wrapf(err, "failed to get user by name")
	}

	sessionID := util.UUID()
	user.SessionID = sessionID

	// create User
	user, err = createUser(ctx, user, a.m, uRepo, uService)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	sRepo := a.sFactory(ctx)
	sService := service.NewSessionService(sRepo, a.m)

	// create Session
	if _, err := createSession(ctx, sessionID, user.ID, a.m, sRepo, sService); err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	return user, nil
}
