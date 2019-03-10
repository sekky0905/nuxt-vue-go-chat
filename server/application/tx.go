package application

import (
	"github.com/pkg/errors"
	"github.com/sekky0905/go-vue-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// CloseTransaction executes after process of tx.
type CloseTransaction func(tx repository.TxManager, err error) error

// beginTxErrorMsg generates and returns tx begin error message.
func beginTxErrorMsg(err error) error {
	return errors.WithStack(&model.SQLError{
		BaseErr:                   err,
		InvalidReasonForDeveloper: model.FailedToBeginTx,
	})
}
