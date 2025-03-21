package database

import (
	"context"
	"database/sql"
)

type Database interface {
	Ping() error
	Connect() error
	Close() error
	QueryRows(query string, args ...any) (*sql.Rows, context.CancelFunc, error)
	QueryRow(query string, args ...any) (*sql.Row, context.CancelFunc, error)
	Exec(query string, args ...any) (sql.Result, context.CancelFunc, error)
}
