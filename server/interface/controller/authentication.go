package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
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
	b, err := GetValueFromPayLoad(r)
	if err != nil {
		ResponseAndLogError(w, err)
		return
	}

	param, err := ParseUserFromPayLoad(b)
	if err != nil {
		ResponseAndLogError(w, err)
		return
	}

	ctx := r.Context()
	user, err := c.aApp.SignUp(ctx, param)
	if err != nil {
		ResponseAndLogError(w, err)
		return
	}

	cookie := c.newCookieWithSessionID(user.SessionID, 86400)
	uDTO := TranslateFromUserToUserDTO(user)

	if err := ResponseWithCookie(w, http.StatusOK, cookie, uDTO); err != nil {
		ResponseAndLogError(w, err)
		return
	}
}

//  ParseUserFromPayLoad parses user from payload.
func ParseUserFromPayLoad(b []byte) (*model.User, error) {
	u := &model.User{}
	if err := json.Unmarshal(b, u); err != nil {
		if err := json.Unmarshal(b, u); err != nil {
			err = &model.InvalidDataError{
				BaseErr:               err,
				DataNameForDeveloper:  "request body",
				DataValueForDeveloper: string(b),
			}
			return nil, errors.WithStack(err)
		}
	}
	return u, nil
}

// newCookieWithSessionID generates and returns cookie with session id.
func (c *authenticationController) newCookieWithSessionID(sessionID string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:   model.SessionIDAtCookie,
		Value:  sessionID,
		MaxAge: maxAge,
	}
}
