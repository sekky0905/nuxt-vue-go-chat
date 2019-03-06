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
