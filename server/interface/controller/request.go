package controller

import (
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// GetValueFromPayLoad load http payload.
func GetValueFromPayLoad(r *http.Request) ([]byte, error) {
	const ContentLength = "Content-Length"
	cl := r.Header.Get(ContentLength)
	length, err := strconv.Atoi(cl)
	if err != nil {
		err = &model.InvalidDataError{
			DataNameForDeveloper:      ContentLength,
			DataValueForDeveloper:     cl,
			InvalidReasonForDeveloper: "content-length should be integer",
		}
		return nil, errors.WithStack(err)
	}

	body := make([]byte, length)
	if _, err = r.Body.Read(body); err != nil && err != io.EOF {
		err = &model.InvalidDataError{
			DataNameForDeveloper:      "request body",
			DataValueForDeveloper:     cl,
			InvalidReasonForDeveloper: "failed to request body",
		}
		return nil, errors.WithStack(err)
	}

	return body, nil
}
