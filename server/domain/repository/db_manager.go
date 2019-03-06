package repository

import (
	"context"
	"database/sql"
)

// SQLManager is the manager of SQL.
type SQLManager interface {
	DBManager
	Beginner
}

// TxManager is the manager of Tx.
type TxManager interface {
	DBManager
	Commit() error
	Rollback() error
}

// DBManager is the manager of DB.
type DBManager interface {
	Querier
	Preparer
	Executor
}

type (
	// Executor is interface of Execute.
	Executor interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	}

	// Preparer is interface of Prepare.
	Preparer interface {
		Prepare(query string) (*sql.Stmt, error)
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	}

	// Querier is interface of Query.
	Querier interface {
		Query(query string, args ...interface{}) (*sql.Rows, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	}

	// Beginner is interface of Begin.
	Beginner interface {
		Begin() (TxManager, error)
	}
)
