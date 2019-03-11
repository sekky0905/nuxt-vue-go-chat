package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	log "github.com/sirupsen/logrus"
)

// Response returns response to client.
func Response(w http.ResponseWriter, statusCode int, obj ...interface{}) error {
	w.Header().Add(ContentLength, "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	var body interface{} = nil
	if len(obj) > 0 {
		body = obj[0]
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		return errors.WithStack(&model.OtherServerError{
			BaseErr:                   err,
			InvalidReasonForDeveloper: "failed to print data to response body",
		})
	}
	return nil
}

// ResponseWithCookie returns response to client with setting of cookie.
func ResponseWithCookie(w http.ResponseWriter, statusCode int, cookie *http.Cookie, obj ...interface{}) error {
	http.SetCookie(w, cookie)

	var body interface{} = nil
	if len(obj) > 0 {
		body = obj[0]
	}

	return Response(w, statusCode, body)
}

// ResponseError returns response error to client.
func ResponseError(w http.ResponseWriter, err error) {
	he := handleError(err)
	if err := Response(w, he.Status, he); err != nil {
		log.Errorf("failed to response:%s", err.Error())
	}
}

// ResponseAndLogError returns response and log error.
func ResponseAndLogError(w http.ResponseWriter, err error) {
	he := handleError(err)
	if he.BaseError != nil {
		log.Errorf("error has occurred:%s, base err = %s", err.Error(), he.BaseError.Error())
	} else {
		log.Errorf("error has occurred:%s", err.Error())
	}

	if err := Response(w, he.Status, he); err != nil {
		log.Errorf("failed to response:%s", err.Error())
	}
}
