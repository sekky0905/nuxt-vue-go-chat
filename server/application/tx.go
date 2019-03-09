package application

import (
	"github.com/pkg/errors"
	"github.com/sekky0905/go-vue-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// CloseTransaction は、Transactionの後処理を行う。
type CloseTransaction func(tx repository.TxManager, err error) error

// beginTxErrorMsg は、TxのBeginエラーメッセージを返す。
func beginTxErrorMsg(err error) error {
	return errors.WithStack(&model.SQLError{
		BaseErr:                   err,
		InvalidReasonForDeveloper: model.FailedToBeginTx,
	})
}
