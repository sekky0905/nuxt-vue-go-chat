package model

import (
	"time"

	"github.com/sekky0905/nuxt-vue-go-chat/server/util"
)

// User is User model.
type User struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	SessionID string    `json:"sessionId"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewUser generates and reruns User.
func NewUser(name, password string) (*User, error) {
	hashed, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:      name,
		Password:  hashed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
