package repository

import (
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// SessionRepository is repository of session.
type SessionRepository interface {
	GetSessionByID(m DBManager, id string) (*model.Session, error)
	InsertSession(m DBManager, user *model.Session) error
	DeleteSession(m DBManager, id string) error
}
