package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// ThreadService is interface of ThreadService.
type ThreadService interface {
	IsAlreadyExistID(ctx context.Context, m repository.SQLManager, id uint32) (bool, error)
	IsAlreadyExistTitle(ctx context.Context, m repository.SQLManager, title string) (bool, error)
}

// threadService is domain service of Thread.
type threadService struct {
	repo repository.ThreadRepository
}

// NewThreadService generates and returns ThreadService.
func NewThreadService(repo repository.ThreadRepository) ThreadService {
	return &threadService{
		repo: repo,
	}
}

// NewThread generates and returns Thread.
func (s threadService) NewThread(title string, user *model.User) *model.Thread {
	return &model.Thread{
		Title: title,
		User:  user,
	}
}

// NewThreadList generates and returns ThreadList.
func (s threadService) NewThreadList(list []*model.Thread, hasNext bool, cursor uint32) *model.ThreadList {
	return &model.ThreadList{
		Threads: list,
		HasNext: hasNext,
		Cursor:  cursor,
	}
}

// IsAlreadyExistID checks duplication of id.
func (s threadService) IsAlreadyExistID(ctx context.Context, m repository.SQLManager, id uint32) (bool, error) {
	var searched *model.Thread
	var err error

	if searched, err = s.repo.GetThreadByID(ctx, m, id); err != nil {
		return false, errors.Wrap(err, "failed to get thread by PropertyNameForDeveloper")
	}
	return searched != nil, nil
}

// IsAlreadyExistID checks duplication of name.
func (s threadService) IsAlreadyExistTitle(ctx context.Context, m repository.SQLManager, title string) (bool, error) {
	var searched *model.Thread
	var err error

	if searched, err = s.repo.GetThreadByTitle(ctx, m, title); err != nil {
		return false, errors.Wrap(err, "failed to get thread by title")
	}
	return searched != nil, nil
}
