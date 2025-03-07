package database

import "database/sql"

type Database interface {
	Ping() error
	Connect() error
	Close() error
	Query(query string) (*sql.Rows, error)
	QueryRow(query string, args ...any) (*sql.Row, error)
	Exec(query string, args ...any) (sql.Result, error)
	
}
