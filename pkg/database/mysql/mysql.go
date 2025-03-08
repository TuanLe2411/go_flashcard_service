package mysql

import (
	"context"
	"database/sql"
	"flashcard_service/pkg/database"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type MySql struct {
	conn    string
	db      *sql.DB
	timeout time.Duration
}

func NewMySql() database.Database {
	config := mysql.Config{
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("MYSQL_URL"),
		DBName:    os.Getenv("MYSQL_DB"),
		Loc:       time.Local,
		ParseTime: true,
	}
	queryTimeout, _ := strconv.Atoi(os.Getenv("MYSQL_QUERY_TIMEOUT_BY_SECOND"))
	return &MySql{
		conn:    config.FormatDSN(),
		timeout: time.Second * time.Duration(queryTimeout),
	}
}

func (m *MySql) Connect() error {
	db, err := sql.Open("mysql", m.conn)
	if err != nil {
		return err
	}
	m.db = db
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetConnMaxIdleTime(time.Minute * 3)
	maxOpenCons, _ := strconv.Atoi(os.Getenv("MYSQL_POOL_MAX_OPEN_CONNECTION"))
	maxIdConst, _ := strconv.Atoi(os.Getenv("MYSQL_POOL_MAX_IDLE_CONNECTION"))
	db.SetMaxIdleConns(maxIdConst)
	db.SetMaxOpenConns(maxOpenCons)
	return m.Ping()
}

func (m *MySql) Close() error {
	return m.db.Close()
}

func (m *MySql) Ping() error {
	return m.db.Ping()
}

func (m *MySql) QueryRows(query string, args ...any) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *MySql) Exec(query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	r, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (m *MySql) QueryRow(query string, args ...any) (*sql.Row, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	row := m.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return row, nil
}
