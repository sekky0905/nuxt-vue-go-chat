package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// RequestManager is the interface of RequestManager.
type RequestManager interface {
	GetValueOfURLParamWithAcceptanceEmpty(r *http.Request, key model.PropertyNameForDeveloper) string
	GetValueOfURLParamWithoutAcceptanceEmpty(r *http.Request, key model.PropertyNameForDeveloper) (string, error)
	GetIntValueOfURLParam(r *http.Request, key model.PropertyNameForDeveloper) (int, error)
	GetUint32ValueOfURLParam(r *http.Request, key model.PropertyNameForDeveloper) (uint32, error)
}

// requestManager is the manager of request.
type requestManager struct {
}

// NewRequestManager generates and returns RequestManager.
func NewRequestManager() RequestManager {
	return &requestManager{}
}

// GetValueOfURLParamWithAcceptanceEmpty gets value specified by key from url.
// This allows empty.
// This is be able to use as follows.
// /hoge/{id}
// /hoge?limit={limit}
func (rm *requestManager) GetValueOfURLParamWithAcceptanceEmpty(r *http.Request, key model.PropertyNameForDeveloper) string {
	if key.String() == "" {
		return ""
	}

	return rm.getValueOfURLParam(r, key)
}

// GetValueOfURLParamWithoutAcceptanceEmpty gets value specified by key from url.
// This doesn't allow empty.
// This is be able to use as follows.
// /hoge/{id}
// /hoge?limit={limit}
func (rm *requestManager) GetValueOfURLParamWithoutAcceptanceEmpty(r *http.Request, key model.PropertyNameForDeveloper) (string, error) {
	if key.String() == "" {
		err := &model.RequiredError{
			PropertyNameForDeveloper: key,
			PropertyNameForUser:      model.PropertyNameKV[key],
		}
		return "", errors.WithStack(err)
	}

	v := rm.getValueOfURLParam(r, key)

	if v == "" {
		err := &model.RequiredError{
			PropertyNameForDeveloper: key,
			PropertyNameForUser:      model.PropertyNameKV[key],
		}
		return "", errors.WithStack(err)
	}

	return v, nil
}

func (rm *requestManager) getValueOfURLParam(r *http.Request, key model.PropertyNameForDeveloper) string {
	return mux.Vars(r)[key.String()]
}

// GetIntValueOfURLParam gets int value specified by key from url.
// This doesn't allow empty.
func (rm *requestManager) GetIntValueOfURLParam(r *http.Request, key model.PropertyNameForDeveloper) (int, error) {
	vStr, err := rm.GetValueOfURLParamWithoutAcceptanceEmpty(r, key)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get value of URL parameter")
	}

	i, err := strconv.Atoi(vStr)
	if err != nil {
		propertyNameForUser := model.PropertyNameKV[key]
		err = &model.InvalidParamError{
			PropertyNameForDeveloper:  key,
			PropertyNameForUser:       propertyNameForUser,
			PropertyValue:             vStr,
			InvalidReasonForDeveloper: fmt.Sprintf("%s should be intger, but requested value is %s", key, vStr),
			InvalidReasonForUser:      fmt.Sprintf("%s は、数字で入力してください", propertyNameForUser),
		}
		return 0, errors.WithStack(err)
	}

	return i, nil
}

// GetUint32ValueOfURLParam gets uint32 value specified by key from url.
// This doesn't allow empty.
func (rm *requestManager) GetUint32ValueOfURLParam(r *http.Request, key model.PropertyNameForDeveloper) (uint32, error) {
	v, err := rm.GetIntValueOfURLParam(r, key)
	if err != nil {
		return model.InvalidID, nil
	}

	return uint32(v), nil
}
