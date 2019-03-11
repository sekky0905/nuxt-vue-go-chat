package controller

import (
	"net/http"

	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/router"
)

// AuthenticationController is the interface of AuthenticationController.
type AuthenticationController interface {
	SignUp(w http.ResponseWriter, r *http.Request)
}

// authenticationController is the controller of authentication.
type authenticationController struct {
	rm   router.RequestManager
	aApp application.AuthenticationService
}

// NewAuthenticationController generates and returns AuthenticationController.
func NewAuthenticationController(rm router.RequestManager, uAPP application.AuthenticationService) AuthenticationController {
	return &authenticationController{
		rm:   rm,
		aApp: uAPP,
	}
}

// SignUp sign up an user.
func (c *authenticationController) SignUp(w http.ResponseWriter, r *http.Request) {
	return
}
