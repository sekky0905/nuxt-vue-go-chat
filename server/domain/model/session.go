package model

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// Session is Session model.
type Session struct {
	ID        string
	UserID    uint32
	CreatedAt time.Time
}

// MarshalLogObject for zap logger.
func (s Session) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", s.ID)
	enc.AddInt32("userID", int32(s.UserID))
	enc.AddTime("createdAt", s.CreatedAt)
	return nil
}
