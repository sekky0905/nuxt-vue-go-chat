package db

import (
	"github.com/pkg/errors"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"
)

// CloseTransaction executes post process of tx.
func CloseTransaction(tx query.TxManager, err error) error {
	if p := recover(); p != nil { // rewrite panic
		err = tx.Rollback()
		err = errors.Wrap(err, "failed to roll back")
		panic(p)
	} else if err != nil {
		err = tx.Rollback()
		err = errors.Wrap(err, "failed to roll back")
	} else {
		err = tx.Commit()
		err = errors.Wrap(err, "failed to commit")
	}
	return err
}
