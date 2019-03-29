package model

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// User is User model.
type User struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	SessionID string    `json:"sessionId"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MarshalLogObject for zap logger.
func (u User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt32("id", int32(u.ID))
	enc.AddString("name", u.Name)
	enc.AddString("sessionID", u.SessionID)
	enc.AddString("password", u.Password)
	enc.AddTime("created_at", u.CreatedAt)
	enc.AddTime("updated_at", u.CreatedAt)
	return nil
}
