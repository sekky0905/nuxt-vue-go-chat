package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"
)

// AuthenticationController is the interface of AuthenticationController.
type AuthenticationController interface {
	InitAuthenticationAPI(g *gin.RouterGroup)
	SignUp(g *gin.Context)
	Login(g *gin.Context)
	Logout(g *gin.Context)
}

// authenticationController is the controller of authentication.
type authenticationController struct {
	aApp application.AuthenticationService
}

// NewAuthenticationController generates and returns AuthenticationController.
func NewAuthenticationController(uAPP application.AuthenticationService) AuthenticationController {
	return &authenticationController{
		aApp: uAPP,
	}
}

// InitAuthenticationAPI initialize Authentication API.
func (c *authenticationController) InitAuthenticationAPI(g *gin.RouterGroup) {
	g.POST("/signUp", c.SignUp)
	g.POST("/login", c.Login)
	g.POST("/logout", c.Logout)
}

// SignUp sign up an user.
func (c *authenticationController) SignUp(g *gin.Context) {
	param := &model.User{}
	if err := g.BindJSON(param); err != nil {
		err = handleValidatorErr(err)
		ResponseAndLogError(g, errors.Wrap(err, "failed to bind json"))
		return
	}

	ctx := g.Request.Context()
	user, err := c.aApp.SignUp(ctx, param)
	if err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to sign up"))
		return
	}

	g.SetCookie(model.SessionIDAtCookie, user.SessionID, 86400, "/", "", true, true)

	uDTO := TranslateFromUserToUserDTO(user)
	g.JSON(http.StatusOK, uDTO)
}

// Login login an user.
func (c *authenticationController) Login(g *gin.Context) {
	param := &model.User{}
	if err := g.BindJSON(param); err != nil {
		err = handleValidatorErr(err)
		ResponseAndLogError(g, errors.Wrap(err, "failed to bind json"))
		return
	}

	ctx := g.Request.Context()
	user, err := c.aApp.Login(ctx, param)
	if err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to sign up"))
		return
	}

	g.SetCookie(model.SessionIDAtCookie, user.SessionID, 86400, "/", "", false, true)

	uDTO := TranslateFromUserToUserDTO(user)
	g.JSON(http.StatusOK, uDTO)
}

// Logout logout an user.
func (c *authenticationController) Logout(g *gin.Context) {
	sessionID, err := g.Cookie(model.SessionIDAtCookie)
	if err != nil && err != http.ErrNoCookie {
		logger.Logger.Warn("failed to read session from Cookie")
		return
	}

	ctx := g.Request.Context()
	if err = c.aApp.Logout(ctx, sessionID); err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to logout"))
		return
	}

	// empty cookie
	g.SetCookie(model.SessionIDAtCookie, "", 0, "", "", false, true)

	g.JSON(http.StatusOK, nil)
}
