package controller

import (
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"
	"gopkg.in/go-playground/validator.v8"
)

// type ValidationErrors map[string]*FieldError

// handleValidatorErr handle validator error.
func handleValidatorErr(err error) error {
	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		logger.Logger.Error("failed to assert ValidationErrors")
	}

	errs := &model.InvalidParamsError{}

	for _, v := range errors {
		e := &model.InvalidParamError{
			BaseErr:       err,
			PropertyName:  model.PropertyName(v.Field),
			PropertyValue: v.Value,
		}

		errs.Errors = append(errs.Errors, e)
	}

	if len(errs.Errors) == 1 {
		return errs.Errors[0]
	}

	return errs
}
