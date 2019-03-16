package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/util"
)

// SessionService is interface of domain service of session.
type SessionService interface {
	NewSession(userID uint32) *model.Session
	SessionID() string
	IsAlreadyExistID(ctx context.Context, id string) (bool, error)
}

// SessionRepoFactory is factory of SessionRepository.
type SessionRepoFactory func(ctx context.Context) repository.SessionRepository

// sessionService is domain service of user.
type sessionService struct {
	m    repository.DBManager
	repo repository.SessionRepository
}

// NewSessionService generates and returns SessionService.
func NewSessionService(m repository.DBManager, repo repository.SessionRepository) SessionService {
	return &sessionService{
		m:    m,
		repo: repo,
	}
}

// NewSession generates and returns Session.
func (s *sessionService) NewSession(userID uint32) *model.Session {
	session := &model.Session{
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	return session
}

// SessionID generates and returns SessionID.
func (s *sessionService) SessionID() string {
	return util.UUID()
}

// IsAlreadyExistID checks whether the data specified by id already exists or not.
func (s sessionService) IsAlreadyExistID(ctx context.Context, id string) (bool, error) {
	var searched *model.Session
	var err error

	if searched, err = s.repo.GetSessionByID(ctx, s.m, id); err != nil {
		return false, errors.Wrap(err, "failed to get session by id")
	}
	return searched != nil, nil
}
