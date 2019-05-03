package controller

import (
	"fmt"
	"net/http"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// handledError is the handled error.
type handledError struct {
	BaseError error   `json:"-"`
	Status    int     `json:"-"`
	Code      ErrCode `json:"code"`
	Message   string  `json:"message"`
}

const systemError = "system error has occurred"

// handleError handles error.
// This generates and returns status code and handledError.
func handleError(err error) *handledError {
	switch errors.Cause(err).(type) {
	case *model.NoSuchDataError:
		realErr, ok := errors.Cause(err).(*model.NoSuchDataError)
		if !ok {
			logger.Logger.Error(fmt.Sprintf("failed to assert. err = %+v", err))
			return nil
		}

		return &handledError{
			BaseError: realErr.BaseErr,
			Status:    http.StatusNotFound,
			Code:      NoSuchDataFailure,
			Message:   errors.Cause(err).Error(),
		}
	case *model.RequiredError:
		realErr, ok := errors.Cause(err).(*model.RequiredError)
		if !ok {
			logger.Logger.Error(fmt.Sprintf("failed to assert. err = %+v", err))
			return nil
		}

		return &handledError{
			BaseError: realErr.BaseErr,
			Status:    http.StatusBadRequest,
			Code:      RequiredFailure,
			Message:   errors.Cause(err).Error(),
		}
	case *model.InvalidParamError:
		realErr, ok := errors.Cause(err).(*model.InvalidParamError)
		if !ok {
			logger.Logger.Error(fmt.Sprintf("failed to assert. err = %+v", err))
			return nil
		}

		return &handledError{
			BaseError: realErr.BaseErr,
			Status:    http.StatusBadRequest,
			Code:      InvalidParameterValueFailure,
			Message:   errors.Cause(err).Error(),
		}
	case *model.InvalidParamsError:
		return &handledError{
			Status:  http.StatusBadRequest,
			Code:    InvalidParametersValueFailure,
			Message: errors.Cause(err).Error(),
		}
	case *model.AlreadyExistError:
		realErr, ok := errors.Cause(err).(*model.AlreadyExistError)
		if !ok {
			logger.Logger.Error(fmt.Sprintf("failed to assert. err = %+v", err))
			return nil
		}

		return &handledError{
			BaseError: realErr.BaseErr,
			Status:    http.StatusConflict,
			Code:      AlreadyExistsFailure,
			Message:   errors.Cause(err).Error(),
		}
	case *model.AuthenticationErr:
		return &handledError{
			Status:  http.StatusUnauthorized,
			Code:    AuthenticationFailure,
			Message: errors.Cause(err).Error(),
		}
	case *model.RepositoryError:
		realErr, ok := errors.Cause(err).(*model.RepositoryError)
		if !ok {
			logger.Logger.Error(fmt.Sprintf("failed to assert. err = %+v", err))
			return nil
		}

		return &handledError{
			BaseError: realErr.BaseErr,
			Status:    http.StatusInternalServerError,
			Code:      InternalDBFailure,
			Message:   systemError,
		}
	case *model.SQLError:
		realErr, ok := errors.Cause(err).(*model.SQLError)
		if !ok {
			logger.Logger.Error(fmt.Sprintf("failed to assert. err = %+v", err))
			return nil
		}

		return &handledError{
			BaseError: realErr.BaseErr,
			Status:    http.StatusInternalServerError,
			Code:      InternalSQLFailure,
			Message:   errors.Cause(err).Error(),
		}
	case *model.OtherServerError:
		realErr, ok := errors.Cause(err).(*model.OtherServerError)
		if !ok {
			logger.Logger.Error(fmt.Sprintf("failed to assert. err = %+v", err))
			return nil
		}

		return &handledError{
			BaseError: realErr.BaseErr,
			Status:    http.StatusInternalServerError,
			Code:      InternalFailure,
			Message:   systemError,
		}
	default:
		return &handledError{
			Status:  http.StatusInternalServerError,
			Code:    InternalFailure,
			Message: systemError,
		}
	}
}
