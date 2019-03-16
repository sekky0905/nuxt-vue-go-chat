package controller

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// handledError is the handled error.
type handledError struct {
	BaseError      error   `json:"-"`
	Status         int     `json:"-"`
	Code           ErrCode `json:"code"`
	Message        string  `json:"message"`
	ErrorUserTitle string  `json:"error_user_title"`
	ErrorUserMsg   string  `json:"error_user_msg"`
}

const systemError = "system error has occurred"
const systemErrorForUser = "システムエラー"

// handleError handles error.
// This generates and returns status code and handledError.
func handleError(err error) *handledError {
	switch errors.Cause(err).(type) {
	case *model.NoSuchDataError:
		realErr := errors.Cause(err).(*model.NoSuchDataError)
		return &handledError{
			BaseError:      realErr.BaseErr,
			Status:         http.StatusNotFound,
			Code:           NoSuchDataFailure,
			Message:        errors.Cause(err).Error(),
			ErrorUserTitle: "不正な指定",
			ErrorUserMsg:   fmt.Sprintf("ご指定された%sのデータが存在しません", realErr.DomainModelNameForUser),
		}
	case *model.RequiredError:
		realErr := errors.Cause(err).(*model.RequiredError)
		return &handledError{
			BaseError:      realErr.BaseErr,
			Status:         http.StatusBadRequest,
			Code:           RequiredFailure,
			Message:        errors.Cause(err).Error(),
			ErrorUserTitle: "入力の不足",
			ErrorUserMsg:   fmt.Sprintf("%sの入力が必要です", realErr.PropertyNameForUser),
		}
	case *model.InvalidParamError:
		realErr := errors.Cause(err).(*model.InvalidParamError)
		return &handledError{
			BaseError:      realErr.BaseErr,
			Status:         http.StatusBadRequest,
			Code:           InvalidParameterValueFailure,
			Message:        errors.Cause(err).Error(),
			ErrorUserTitle: "不正な入力",
			ErrorUserMsg:   realErr.InvalidReasonForUser,
		}
	case *model.InvalidParamsError:
		return &handledError{
			Status:         http.StatusBadRequest,
			Code:           InvalidParametersValueFailure,
			Message:        errors.Cause(err).Error(),
			ErrorUserTitle: "不正な入力",
			ErrorUserMsg:   "不正な入力です",
		}
	case *model.AlreadyExistError:
		realErr := errors.Cause(err).(*model.AlreadyExistError)
		return &handledError{
			BaseError:      realErr.BaseErr,
			Status:         http.StatusConflict,
			Code:           AlreadyExistsFailure,
			Message:        errors.Cause(err).Error(),
			ErrorUserTitle: "不正な入力",
			ErrorUserMsg:   fmt.Sprintf("ご指定いただいた%sのデータは既に存在しています", realErr.DomainModelNameForUser),
		}
	case *model.AuthenticationErr:
		return &handledError{
			Status:         http.StatusUnauthorized,
			Code:           AuthenticationFailure,
			Message:        errors.Cause(err).Error(),
			ErrorUserTitle: "認証エラー",
			ErrorUserMsg:   "認証に失敗しました、IDもしくはパスワードが不正か既に利用されています",
		}
	case *model.RepositoryError:
		realErr := errors.Cause(err).(*model.RepositoryError)
		return &handledError{
			BaseError:      realErr.BaseErr,
			Status:         http.StatusInternalServerError,
			Code:           InternalDBFailure,
			Message:        systemError,
			ErrorUserTitle: systemErrorForUser,
			ErrorUserMsg:   fmt.Sprintf("[エラーコード: %s]システムエラーが発生しました。", InternalDBFailure),
		}
	case *model.SQLError:
		realErr := errors.Cause(err).(*model.SQLError)
		return &handledError{
			BaseError:      realErr.BaseErr,
			Status:         http.StatusInternalServerError,
			Code:           InternalSQLFailure,
			Message:        errors.Cause(err).Error(),
			ErrorUserTitle: systemErrorForUser,
			ErrorUserMsg:   fmt.Sprintf("[エラーコード: %s]システムエラーが発生しました。", InternalSQLFailure),
		}
	case *model.OtherServerError:
		realErr := errors.Cause(err).(*model.OtherServerError)
		return &handledError{
			BaseError:      realErr.BaseErr,
			Status:         http.StatusInternalServerError,
			Code:           InternalFailure,
			Message:        systemError,
			ErrorUserTitle: systemErrorForUser,
			ErrorUserMsg:   fmt.Sprintf("[エラーコード: %s]システムエラーが発生しました。", ServerError),
		}
	default:
		return &handledError{
			Status:         http.StatusInternalServerError,
			Code:           InternalFailure,
			Message:        systemError,
			ErrorUserTitle: systemErrorForUser,
			ErrorUserMsg:   "システムエラーが発生しました。",
		}
	}
}
