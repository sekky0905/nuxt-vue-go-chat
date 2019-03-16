package controller

import (
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	log "github.com/sirupsen/logrus"

	"gopkg.in/go-playground/validator.v8"
)

// type ValidationErrors map[string]*FieldError

// handleValidatorErr handle validator error.
func handleValidatorErr(err error) error {
	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		log.Warnf("failed to assert ValidationErrors")
	}

	errs := &model.InvalidParamsError{}

	for _, v := range errors {
		e := &model.InvalidParamError{
			BaseErr:                  err,
			PropertyNameForDeveloper: model.PropertyNameForDeveloper(v.Field),
			PropertyNameForUser:      model.PropertyNameKV[model.PropertyNameForDeveloper(v.Field)],
			PropertyValue:            v.Value,
		}

		errs.Errors = append(errs.Errors, e)
	}

	if len(errs.Errors) == 1 {
		return errs.Errors[0]
	}

	return errs
}
