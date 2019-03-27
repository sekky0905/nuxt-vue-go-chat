package application

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
)

// ThreadService is interface of ThreadService.
type ThreadService interface {
	ListThreads(ctx context.Context, limit int, cursor uint32) (*model.ThreadList, error)
	GetThread(ctx context.Context, id uint32) (*model.Thread, error)
	CreateThread(ctx context.Context, thread *model.Thread) (*model.Thread, error)
	UpdateThread(ctx context.Context, id uint32, thread *model.Thread) (*model.Thread, error)
	DeleteThread(ctx context.Context, id uint32) error
}

// threadService is application service of thread.
type threadService struct {
	m        repository.DBManager
	service  service.ThreadService
	repo     repository.ThreadRepository
	txCloser CloseTransaction
}

// NewThreadService generates and returns ThreadService.
func NewThreadService(m repository.DBManager, service service.ThreadService, repo repository.ThreadRepository, txCloser CloseTransaction) ThreadService {
	return &threadService{
		m:        m,
		service:  service,
		repo:     repo,
		txCloser: txCloser,
	}
}

// ListThreads gets ThreadList.
func (a *threadService) ListThreads(ctx context.Context, limit int, cursor uint32) (*model.ThreadList, error) {
	threads, err := a.repo.ListThreads(ctx, a.m, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to thread threads")
	}

	return threads, nil
}

// GetThread gets Thread.
func (a *threadService) GetThread(ctx context.Context, id uint32) (*model.Thread, error) {
	thread, err := a.repo.GetThreadByID(ctx, a.m, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get thread by PropertyNameForDeveloper")
	}

	return thread, nil
}

// CreateThread creates Thread.
func (a *threadService) CreateThread(ctx context.Context, param *model.Thread) (thread *model.Thread, err error) {
	tx, err := a.m.Begin()
	if err != nil {
		return nil, beginTxErrorMsg(err)
	}

	defer func() {
		if err := a.txCloser(tx, err); err != nil {
			err = errors.Wrap(err, "failed to close tx")
		}
	}()

	param.CreatedAt = time.Now()
	param.UpdatedAt = time.Now()

	yes, err := a.service.IsAlreadyExistTitle(ctx, param.Title)
	if yes {
		err = &model.AlreadyExistError{
			PropertyNameForDeveloper:    model.TitlePropertyForDeveloper,
			PropertyNameForUser:         model.TitlePropertyForUser,
			PropertyValue:               param.Title,
			DomainModelNameForDeveloper: model.DomainModelNameThreadForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameThreadForUser,
		}
		return nil, errors.Wrap(err, "already exist title")
	}

	if _, ok := errors.Cause(err).(*model.NoSuchDataError); !ok {
		return nil, errors.Wrap(err, "failed is already exist title")
	}

	id, err := a.repo.InsertThread(ctx, tx, param)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert thread")
	}
	param.ID = id
	return param, nil
}

// UpdateThread updates Thread.
func (a *threadService) UpdateThread(ctx context.Context, id uint32, param *model.Thread) (thread *model.Thread, err error) {
	tx, err := a.m.Begin()
	if err != nil {
		return nil, beginTxErrorMsg(err)
	}

	defer func() {
		if err := a.txCloser(tx, err); err != nil {
			err = errors.Wrap(err, "failed to close tx")
		}
	}()

	yes, err := a.service.IsAlreadyExistID(ctx, param.ID)
	if !yes {
		err = &model.NoSuchDataError{
			PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
			PropertyNameForUser:         model.IDPropertyForUser,
			PropertyValue:               param.ID,
			DomainModelNameForDeveloper: model.DomainModelNameThreadForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameThreadForUser,
		}
		return nil, errors.Wrap(err, "does not exists ID")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to is already exist ID")
	}

	param.UpdatedAt = time.Now()

	if err := a.repo.UpdateThread(ctx, tx, param.ID, param); err != nil {
		return nil, errors.Wrap(err, "failed to update thread")
	}

	return param, nil
}

// DeleteThread deletes Thread.
func (a *threadService) DeleteThread(ctx context.Context, id uint32) (err error) {
	tx, err := a.m.Begin()
	if err != nil {
		return beginTxErrorMsg(err)
	}

	defer func() {
		if err := a.txCloser(tx, err); err != nil {
			err = errors.Wrap(err, "failed to close tx")
		}
	}()

	yes, err := a.service.IsAlreadyExistID(ctx, id)
	if !yes {
		err = &model.NoSuchDataError{
			PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
			PropertyNameForUser:         model.IDPropertyForUser,
			PropertyValue:               id,
			DomainModelNameForDeveloper: model.DomainModelNameThreadForDeveloper,
			DomainModelNameForUser:      model.DomainModelNameThreadForUser,
		}
		return errors.Wrap(err, "does not exists PropertyNameForDeveloper")
	}
	if err != nil {
		return errors.Wrap(err, "failed to is already exist PropertyNameForDeveloper")
	}

	if err := a.repo.DeleteThread(ctx, tx, id); err != nil {
		return errors.Wrap(err, "failed to delete thread")
	}

	return nil
}
