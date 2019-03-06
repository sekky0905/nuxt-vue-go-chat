package model

import (
	"time"
)

// Session is Session model.
type Session struct {
	ID        string
	UserID    uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewSession generates and returns Session.
func NewSession(userID uint32) *Session {
	session := &Session{
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return session
}
