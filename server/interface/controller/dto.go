package controller

import (
	"time"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"
	"go.uber.org/zap"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// UserDTO is DTO of User.
type UserDTO struct {
	ID        uint32    `json:"id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	SessionID string    `json:"sessionId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TranslateFromUserToUserDTO translates from User to UserDTO.
func TranslateFromUserToUserDTO(user *model.User) *UserDTO {
	return &UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		SessionID: user.SessionID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ThreadDTO is DTO of Thread.
type ThreadDTO struct {
	ID        uint32 `json:"id"`
	Title     string `json:"title" binding:"required"`
	*UserDTO  `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TranslateFromThreadDTOToThread translates from ThreadDTO to Thread.
func TranslateFromThreadDTOToThread(dto *ThreadDTO) *model.Thread {
	return &model.Thread{
		ID:    dto.ID,
		Title: dto.Title,
		User: &model.User{
			ID:        dto.UserDTO.ID,
			Name:      dto.UserDTO.Name,
			CreatedAt: dto.UserDTO.CreatedAt,
			UpdatedAt: dto.UserDTO.UpdatedAt,
		},
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

// CommentDTO is DTO of CommentD.
type CommentDTO struct {
	ID        uint32 `json:"id"`
	Content   string `json:"content" binding:"required"`
	ThreadID  uint32 `json:"threadId" binding:"required"`
	*UserDTO  `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TranslateFromCommentDTOToComment translates from CommentDTO to Comment.
func TranslateFromCommentDTOToComment(dto *CommentDTO) *model.Comment {
	logger.Logger.Info("CO", zap.String("content", dto.Content))
	return &model.Comment{
		ID:      dto.ID,
		Content: dto.Content,
		User: &model.User{
			ID:        dto.UserDTO.ID,
			Name:      dto.UserDTO.Name,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
		ThreadID:  dto.ThreadID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
