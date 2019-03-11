package application

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
)

// AuthenticationService is the interface of AuthenticationService.
type AuthenticationService interface {
	SignUp(ctx context.Context, param *model.User) (*model.User, error)
}

// AuthenticationServiceDIInput is DI input of AuthenticationService.
type AuthenticationServiceDIInput struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	userService       service.UserService
	sessionService    service.SessionService
}

// NewAuthenticationServiceDIInput generates and returns AuthenticationServiceDIInput.
func NewAuthenticationServiceDIInput(uRepo repository.UserRepository, sRepo repository.SessionRepository, uService service.UserService, sService service.SessionService) *AuthenticationServiceDIInput {
	return &AuthenticationServiceDIInput{
		userRepository:    uRepo,
		sessionRepository: sRepo,
		userService:       uService,
		sessionService:    sService,
	}
}

// authenticationService is the service of authentication.
type authenticationService struct {
	m                 repository.DBManager
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	userService       service.UserService
	sessionService    service.SessionService
	txCloser          CloseTransaction
}

// NewAuthenticationService generates and returns AuthenticationService.
func NewAuthenticationService(m repository.DBManager, diInput AuthenticationServiceDIInput, txCloser CloseTransaction) AuthenticationService {
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

	sessionID := s.sessionService.SessionID()
	user.SessionID = sessionID

	// create User
	user, err = s.createUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	// create Session
	if _, err := s.createSession(ctx, sessionID, user.ID); err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	return user, nil
}

// createUser creates the user.
func (s *authenticationService) createUser(ctx context.Context, param *model.User) (*model.User, error) {
	// not allow duplicated name.
	yes, err := s.userService.IsAlreadyExistName(ctx, param.Name)
	if yes {
		err = &model.AlreadyExistError{
			PropertyNameForDeveloper:    model.NamePropertyForDeveloper,
			PropertyNameForUser:         model.NamePropertyForUser,
			PropertyValue:               param.Name,
			DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameUserForUser,
		}

		return nil, errors.Wrap(err, "failed to check whether already exists name or not")
	}

	if err != nil {
		if _, ok := errors.Cause(err).(*model.NoSuchDataError); !ok {
			return nil, errors.Wrap(err, "failed to check whether already exists name or not")
		}
	}

	id, err := s.userRepository.InsertUser(s.m, param)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert user")
	}
	param.ID = id

	return param, nil
}

// createSession creates the session.
func (s *authenticationService) createSession(ctx context.Context, sessionID string, userID uint32) (*model.Session, error) {
	session := model.NewSession(userID)
	session.ID = sessionID

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

	if err := s.sessionRepository.InsertSession(s.m, session); err != nil {
		return nil, errors.Wrap(err, "failed to insert session")
	}
	return session, nil
}
