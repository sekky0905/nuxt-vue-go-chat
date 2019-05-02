package db

import (
	"context"
	"database/sql"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"

	// SQL Driver.
	_ "github.com/go-sql-driver/mysql"
)

// dbManager is the manager of SQL.
type dbManager struct {
	Conn *sql.DB
}

// NewDBManager generates and returns DBManager.
func NewDBManager() query.DBManager {
	conn, err := sql.Open("mysql", "root:@tcp(nvgdb:3306)/nuxt_vue_go_chat?charset=utf8mb4&parseTime=True")
	if err != nil {
		panic(err.Error())
	}

	return &dbManager{
		Conn: conn,
	}
}

// Exec executes SQL.
func (s dbManager) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.Conn.Exec(query, args...)
}

// ExecContext executes SQL with context.
func (s *dbManager) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.Conn.ExecContext(ctx, query, args...)
}

// Query executes query which return row.
func (s *dbManager) Query(query string, args ...interface{}) (*sql.Rows, error) {
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
func (s *dbManager) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
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
func (s *dbManager) Prepare(query string) (*sql.Stmt, error) {
	return s.Conn.Prepare(query)
}

// Prepare prepares statement for Query and Exec later with context.
func (s *dbManager) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.Conn.PrepareContext(ctx, query)
}

// Begin begins tx.
func (s *dbManager) Begin() (query.TxManager, error) {
	return s.Conn.Begin()
}
