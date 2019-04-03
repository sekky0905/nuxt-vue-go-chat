package controller

import (
	"net/http"
	"strconv"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"
	"go.uber.org/zap"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"

	"github.com/pkg/errors"

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
	g.GET("/threads/:threadId/comments", c.ListComments)
	g.GET("/threads/:threadId/comments/:id", c.GetComment)
	g.POST("/threads/:threadId/comments", c.CreateComment)
	g.PUT("/threads/:threadId/comments/:id", c.UpdateComment)
	g.DELETE("/threads/:threadId/comments/:id", c.DeleteComment)
}

// NewCommentController generates and returns CommentController.
func NewCommentController(cAPP application.CommentService) CommentController {
	return &commentController{
		cApp: cAPP,
	}
}

// ListComment gets CommentList.
func (c commentController) ListComments(g *gin.Context) {
	limit, err := strconv.Atoi(g.Query("limit"))
	if err != nil {
		limit = defaultLimit
	}

	cursorInt, err := strconv.Atoi(g.Query("cursor"))
	if err != nil {
		cursorInt = defaultCursor
	}

	cursor := uint32(cursorInt)

	threadIInt, err := strconv.Atoi(g.Param("threadId"))
	if err != nil || threadIInt < 1 {
		err = &model.InvalidParamError{
			BaseErr:                   err,
			PropertyNameForDeveloper:  model.ThreadIDPropertyForDeveloper,
			InvalidReasonForDeveloper: "threadId should be number and over 0",
		}

		ResponseAndLogError(g, errors.Wrap(err, "failed to list comments"))
		return
	}

	threadID := uint32(threadIInt)

	ctx := g.Request.Context()
	comment, err := c.cApp.ListComments(ctx, threadID, limit, cursor)
	if err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to list comments"))
		return
	}

	g.JSON(http.StatusOK, comment)

}

// GetComment gets Comment.
func (c commentController) GetComment(g *gin.Context) {
	idInt, err := strconv.Atoi(g.Param("id"))
	if err != nil || idInt < 1 {
		logger.Logger.Info("DEBUG>>>", zap.String("ERR", err.Error()))
		err = &model.InvalidParamError{
			BaseErr:                   err,
			PropertyNameForDeveloper:  model.IDPropertyForDeveloper,
			InvalidReasonForDeveloper: "id should be number and over 0",
		}

		logger.Logger.Info("DEBUG", zap.String("ERR", err.Error()))

		ResponseAndLogError(g, err)
		return
	}

	id := uint32(idInt)

	logger.Logger.Info("DEBUG", zap.Int("ID", idInt))

	ctx := g.Request.Context()
	comment, err := c.cApp.GetComment(ctx, id)
	if err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to get comment"))
		return
	}
	g.JSON(http.StatusOK, comment)
}

// CreateComment creates Comment.
func (c commentController) CreateComment(g *gin.Context) {
	dto := &CommentDTO{}
	if err := g.BindJSON(dto); err != nil {
		err = handleValidatorErr(err)
		ResponseAndLogError(g, errors.Wrap(err, "failed to bind json"))
		return
	}

	param := TranslateFromCommentDTOToComment(dto)

	ctx := g.Request.Context()
	thread, err := c.cApp.CreateComment(ctx, param)
	if err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed created comment"))
		return
	}

	g.JSON(http.StatusOK, thread)
}

// UpdateComment updates Comment.
func (c commentController) UpdateComment(g *gin.Context) {
	dto := &CommentDTO{}
	if err := g.BindJSON(dto); err != nil {
		err = handleValidatorErr(err)
		ResponseAndLogError(g, errors.Wrap(err, "failed to bind json"))
		return
	}

	idInt, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		err = &model.InvalidParamError{
			BaseErr:                  err,
			PropertyNameForDeveloper: model.IDPropertyForDeveloper,
			PropertyValue:            g.Param("id"),
		}
		err = handleValidatorErr(err)
		ResponseAndLogError(g, errors.Wrap(err, "failed to change id from string to int"))
		return
	}

	id := uint32(idInt)

	param := TranslateFromCommentDTOToComment(dto)

	ctx := g.Request.Context()
	thread, err := c.cApp.UpdateComment(ctx, id, param)
	if err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to update comment"))
		return
	}

	g.JSON(http.StatusOK, thread)

}

// DeleteComment deletes Comment.
func (c commentController) DeleteComment(g *gin.Context) {
	idInt, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		err = &model.InvalidParamError{
			BaseErr:                  err,
			PropertyNameForDeveloper: model.IDPropertyForDeveloper,
			PropertyValue:            g.Param("id"),
		}
		err = handleValidatorErr(err)
		ResponseAndLogError(g, errors.Wrap(err, "failed to change id from string to int"))
		return
	}

	id := uint32(idInt)

	ctx := g.Request.Context()
	if err := c.cApp.DeleteComment(ctx, id); err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to delete comment"))
		return
	}

	g.JSON(http.StatusOK, nil)
}
