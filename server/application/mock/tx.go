package mock_application

import "github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"

// MockCloseTransaction executes after process of tx.
func MockCloseTransaction(tx repository.TxManager, err error) error {
	return nil
}
