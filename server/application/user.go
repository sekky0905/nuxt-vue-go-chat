package application

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
)

// UserService is the interface of UserService.
type UserService interface {
}

// userService is the service of user.
type userService struct {
	m        repository.SQLManager
	uFactory service.UserRepoFactory
	txCloser CloseTransaction
}

// NewUserService は、generates and returns UserService.
func NewUserService(m repository.SQLManager, f service.UserRepoFactory, txCloser CloseTransaction) UserService {
	return &userService{
		m:        m,
		uFactory: f,
		txCloser: txCloser,
	}
}

// createUser creates the user.
func createUser(ctx context.Context, param *model.User, db repository.DBManager, repo repository.UserRepository, uService service.UserService) (*model.User, error) {
	param.CreatedAt = time.Now()
	param.UpdatedAt = time.Now()

	// not allow duplicated name.
	yes, err := uService.IsAlreadyExistName(ctx, param.Name)
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

	if _, ok := errors.Cause(err).(*model.NoSuchDataError); !ok {
		return nil, errors.Wrap(err, "failed to check whether already exists name or not")
	}

	id, err := repo.InsertUser(db, param)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert user")
	}
	param.ID = id

	return param, nil
}
