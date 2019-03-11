package controller

import (
	"time"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// UserDTO is DTO of User.
type UserDTO struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	SessionID string    `json:"sessionId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TranslateFromUserToUserDTO translate from User to UserDTO.
func TranslateFromUserToUserDTO(user *model.User) *UserDTO {
	return &UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		SessionID: user.SessionID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
