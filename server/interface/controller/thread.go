package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// ThreadController is the interface of ThreadController.
type ThreadController interface {
	InitThreadAPI(g *gin.RouterGroup)
	ListThreads(g *gin.Context)
	GetThread(g *gin.Context)
	CreateThread(g *gin.Context)
	UpdateThread(g *gin.Context)
	DeleteThread(g *gin.Context)
}

// threadController is the controller of thread.
type threadController struct {
	tApp application.ThreadService
}

// NewThreadController generates and returns ThreadController.
func NewThreadController(tApp application.ThreadService) ThreadController {
	return &threadController{
		tApp: tApp,
	}
}

// InitThreadAPI initialize Thread API.
func (c *threadController) InitThreadAPI(g *gin.RouterGroup) {
	g.GET("/threads", c.ListThreads)
	g.GET("/threads/:id", c.GetThread)
	g.POST("/threads", c.CreateThread)
	g.PUT("/threads/:id", c.UpdateThread)
	g.DELETE("/threads/:id", c.DeleteThread)
}

// ListThreads gets ThreadList.
func (c *threadController) ListThreads(g *gin.Context) {
	limit, err := strconv.Atoi(g.Query("limit"))
	if err != nil {
		limit = defaultLimit
	}

	cursorInt, err := strconv.Atoi(g.Query("cursor"))
	if err != nil {
		cursorInt = defaultCursor
	}

	cursor := uint32(cursorInt)

	ctx := g.Request.Context()
	thread, err := c.tApp.ListThreads(ctx, limit, cursor)
	if err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to sign up"))
		return
	}

	g.JSON(http.StatusOK, thread)
}

// GetThread gets Thread.
func (c *threadController) GetThread(g *gin.Context) {
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
	thread, err := c.tApp.GetThread(ctx, id)
	if err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to sign up"))
		return
	}

	g.JSON(http.StatusOK, thread)
}
