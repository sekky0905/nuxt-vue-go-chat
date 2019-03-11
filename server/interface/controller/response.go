package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sekky0905/go-vue-chat/server/domain/model"
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

// TODO implement of error response.
