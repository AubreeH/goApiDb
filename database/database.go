package database

import (
	"context"
	"database/sql"
	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq"
	// _ "modernc.org/sqlite"
)

type DbInstance interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	getTableColumns() map[string][]string
}

// Database contains the [sql.DB] database connection and basic database info.
type Database struct {
	Db           *sql.DB
	dbName       string
	tableColumns map[string][]string
}
