package model

import (
	"time"
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
