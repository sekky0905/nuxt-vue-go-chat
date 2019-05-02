package repository

import (
	"context"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// SessionRepository is repository of session.
type SessionRepository interface {
	GetSessionByID(ctx context.Context, m query.SQLManager, id string) (*model.Session, error)
	InsertSession(ctx context.Context, m query.SQLManager, session *model.Session) error
	DeleteSession(ctx context.Context, m query.SQLManager, id string) error
}
