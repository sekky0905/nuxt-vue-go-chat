package application

import (
	"context"
	"time"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
)

// AuthenticationService is the interface of AuthenticationService.
type AuthenticationService interface {
	SignUp(ctx context.Context, param *model.User) (*model.User, error)
	Login(ctx context.Context, param *model.User) (*model.User, error)
}

// AuthenticationServiceDIInput is DI input of AuthenticationService.
type AuthenticationServiceDIInput struct {
	userRepository        repository.UserRepository
	sessionRepository     repository.SessionRepository
	userService           service.UserService
	sessionService        service.SessionService
	authenticationService service.AuthenticationService
}

// NewAuthenticationServiceDIInput generates and returns AuthenticationServiceDIInput.
func NewAuthenticationServiceDIInput(uRepo repository.UserRepository, sRepo repository.SessionRepository, uService service.UserService, sService service.SessionService, aService service.AuthenticationService) *AuthenticationServiceDIInput {
	return &AuthenticationServiceDIInput{
		userRepository:        uRepo,
		sessionRepository:     sRepo,
		userService:           uService,
		sessionService:        sService,
		authenticationService: aService,
	}
}

// authenticationService is the service of authentication.
type authenticationService struct {
	m                     repository.DBManager
	userRepository        repository.UserRepository
	sessionRepository     repository.SessionRepository
	userService           service.UserService
	sessionService        service.SessionService
	authenticationService service.AuthenticationService
	txCloser              CloseTransaction
}

// NewAuthenticationService generates and returns AuthenticationService.
func NewAuthenticationService(m repository.DBManager, diInput *AuthenticationServiceDIInput, txCloser CloseTransaction) AuthenticationService {
	return &authenticationService{
		m:                 m,
		userRepository:    diInput.userRepository,
		sessionRepository: diInput.sessionRepository,
		userService:       diInput.userService,
		sessionService:    diInput.sessionService,
		txCloser:          txCloser,
	}
}

// SignUp sign up an user.
func (s *authenticationService) SignUp(ctx context.Context, param *model.User) (user *model.User, err error) {
	tx, err := s.m.Begin()
	if err != nil {
		return nil, beginTxErrorMsg(err)
	}

	defer func() {
		if err := s.txCloser(tx, err); err != nil {
			err = errors.Wrap(err, "failed to close tx")
		}
	}()

	user, err = s.userService.NewUser(param.Name, param.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new user")
	}

	sessionID := s.sessionService.SessionID()
	user.SessionID = sessionID

	// create User
	user, err = s.createUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	session := s.sessionService.NewSession(user.ID)
	session.ID = user.SessionID

	// create Session
	if _, err := s.createSession(ctx, session); err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	return user, nil
}

// createUser creates the user.
func (s *authenticationService) createUser(ctx context.Context, user *model.User) (*model.User, error) {
	// not allow duplicated name.
	yes, err := s.userService.IsAlreadyExistName(ctx, user.Name)
	if yes {
		err = &model.AlreadyExistError{
			PropertyNameForDeveloper:    model.NamePropertyForDeveloper,
			PropertyNameForUser:         model.NamePropertyForUser,
			PropertyValue:               user.Name,
			DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameUserForUser,
		}

		return nil, errors.WithStack(err)
	}

	if err != nil {
		if _, ok := errors.Cause(err).(*model.NoSuchDataError); !ok {
			return nil, errors.Wrap(err, "failed to check whether already exists name or not")
		}
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	logger.Logger.Info("XXXXXX", zap.Object("user", user))

	id, err := s.userRepository.InsertUser(ctx, s.m, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert user")
	}
	user.ID = id

	return user, nil
}

// createSession creates the session.
func (s *authenticationService) createSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	// ready for collision of UUID.
	yes := true
	var err error
	for yes {
		yes, err = s.sessionService.IsAlreadyExistID(ctx, session.ID)
		if err != nil {
			if _, ok := errors.Cause(err).(*model.NoSuchDataError); !ok {
				return nil, errors.Wrap(err, "failed to check whether already exists id or not")
			}
		}
	}

	session.CreatedAt = time.Now()

	if err := s.sessionRepository.InsertSession(ctx, s.m, session); err != nil {
		return nil, errors.Wrap(err, "failed to insert session")
	}
	return session, nil
}

// Login Login an user.
func (s *authenticationService) Login(ctx context.Context, param *model.User) (user *model.User, err error) {
	tx, err := s.m.Begin()
	if err != nil {
		return nil, beginTxErrorMsg(err)
	}

	defer func() {
		if err := s.txCloser(tx, err); err != nil {
			err = errors.Wrap(err, "failed to close tx:%#v")
		}
	}()

	logger.Logger.Info("DDDD", zap.Object("param", param))

	ok, user, err := s.authenticationService.Authenticate(ctx, param.Name, param.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to authenticate")
	} else if !ok {
		return nil, errors.WithStack(&model.AuthenticationErr{
			BaseErr: errors.New("name or pass is invalid"),
		})
	}

	session := s.sessionService.NewSession(user.ID)
	session.ID = s.sessionService.SessionID()

	session, err = s.createSession(ctx, session)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	user.SessionID = session.ID

	if err := s.userRepository.UpdateUser(ctx, s.m, user.ID, user); err != nil {
		return nil, errors.Wrap(err, "failed to insert user")
	}

	return user, nil
}
