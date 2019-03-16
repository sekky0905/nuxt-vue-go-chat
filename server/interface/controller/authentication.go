package controller

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// AuthenticationController is the interface of AuthenticationController.
type AuthenticationController interface {
	InitAuthenticationAPI(g *gin.RouterGroup)
	SignUp(g *gin.Context)
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
}

// SignUp sign up an user.
func (c *authenticationController) SignUp(g *gin.Context) {
	param := &model.User{}
	if err := g.BindJSON(param); err != nil {
		err = handleValidatorErr(err)
		logrus.Infof("EEEEEEE===%#v", err)
		ResponseAndLogError(g, errors.Wrap(err, "failed to bind json"))
		return
	}

	logrus.Info("======mksvegijorwkpeql@go2")

	ctx := g.Request.Context()
	user, err := c.aApp.SignUp(ctx, param)
	if err != nil {
		ResponseAndLogError(g, errors.Wrap(err, "failed to sign up"))
		return
	}

	g.SetCookie(model.SessionIDAtCookie, user.SessionID, 86400, "", "", true, true)

	uDTO := TranslateFromUserToUserDTO(user)
	g.JSON(http.StatusOK, uDTO)
}
