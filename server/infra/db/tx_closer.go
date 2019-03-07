package db

import (
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// CloseTransaction は、Transactionの後処理を行う。
func CloseTransaction(tx repository.TxManager, err error) error {
	if p := recover(); p != nil {
		err = tx.Rollback()
		panic(p) // panicをもう一回
	} else if err != nil {
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}
	return errors.WithStack(err)
}
