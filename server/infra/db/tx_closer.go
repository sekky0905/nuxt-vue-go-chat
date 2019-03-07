package db

import (
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
)

// CloseTransaction executes post process of tx.
func CloseTransaction(tx repository.TxManager, err error) error {
	if p := recover(); p != nil { // rewrite panic
		err = tx.Rollback()
		panic(p)
	} else if err != nil {
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}
	return errors.WithStack(err)
}
