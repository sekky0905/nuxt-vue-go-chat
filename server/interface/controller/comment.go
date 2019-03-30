package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
)

// CommentController is the interface of CommentController.
type CommentController interface {
	ListComments(g *gin.Context)
	GetComment(g *gin.Context)
	CreateComment(g *gin.Context)
	UpdateComment(g *gin.Context)
	DeleteComment(g *gin.Context)
}

// commentController is the controller of comment.
type commentController struct {
	cApp application.CommentService
}

// InitCommentAPI initialize Comment API.
func (c *commentController) InitCommentAPI(g *gin.RouterGroup) {
	g.GET("/comments", c.ListComments)
	g.GET("/comments/:id", c.GetComment)
	g.POST("/comments", c.CreateComment)
	g.PUT("/comments/:id", c.UpdateComment)
	g.DELETE("/comments/:id", c.DeleteComment)
}

// NewCommentController generates and returns CommentController.
func NewCommentController(cAPP application.CommentService) CommentController {
	return &commentController{
		cApp: cAPP,
	}
}

// ListComment gets CommentList.
func (c commentController) ListComments(g *gin.Context) {

}

// GetComment gets Comment.
func (c commentController) GetComment(g *gin.Context) {
}

// CreateComment creates Comment.
func (c commentController) CreateComment(g *gin.Context) {
}

// UpdateComment updates Comment.
func (c commentController) UpdateComment(g *gin.Context) {

}

// DeleteComment deletes Comment.
func (c commentController) DeleteComment(g *gin.Context) {
}
