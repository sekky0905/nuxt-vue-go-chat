package repository

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// SessionRepository is repository of session.
type SessionRepository interface {
	GetSessionByID(m SQLManager, id string) (*model.Session, error)
	InsertSession(m SQLManager, user *model.Session) error
	DeleteSession(m SQLManager, id string) error
}

// SessionRepoFactory is Factory of SessionRepository.
type SessionRepoFactory func(ctx context.Context) SessionRepository
