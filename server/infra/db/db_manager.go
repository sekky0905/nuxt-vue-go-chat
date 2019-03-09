package db

import (
	"context"
	"database/sql"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"

	// SQL Driver.
	_ "github.com/go-sql-driver/mysql"
)

// sqlManager is the manager of SQL.
type sqlManager struct {
	Conn *sql.DB
}

// NewSQLManager generates and returns SQLManager.
func NewSQLManager() repository.SQLManager {
	conn, err := sql.Open("mysql", "root:@tcp(nvgdb:3306)/nuxt_vue_go_chat?charset=utf8mb4&parseTime=True")
	if err != nil {
		panic(err.Error())
	}

	return &sqlManager{
		Conn: conn,
	}
}

// Exec executes SQL.
func (s sqlManager) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.Conn.Exec(query, args...)
}

// ExecContext executes SQL with context.
func (s *sqlManager) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.Conn.ExecContext(ctx, query, args...)
}

// Query executes query which return row.
func (s *sqlManager) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.Conn.Query(query, args...)
	if err != nil {
		err = &model.SQLError{
			BaseErr:                   err,
			InvalidReasonForDeveloper: "failed to execute query",
		}
		return nil, err
	}

	return rows, nil
}

// QueryContext executes query which return row with context.
func (s *sqlManager) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.Conn.Query(query, args...)
	if err != nil {
		err = &model.SQLError{
			BaseErr:                   err,
			InvalidReasonForDeveloper: "failed to execute query with context",
		}
		return nil, err
	}
	return rows, nil
}

// Prepare prepares statement for Query and Exec later.
func (s *sqlManager) Prepare(query string) (*sql.Stmt, error) {
	return s.Conn.Prepare(query)
}

// Prepare prepares statement for Query and Exec later with context.
func (s *sqlManager) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.Conn.PrepareContext(ctx, query)
}

// Begin begins tx.
func (s *sqlManager) Begin() (repository.TxManager, error) {
	return s.Conn.Begin()
}
