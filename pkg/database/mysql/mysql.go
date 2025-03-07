package mysql

import (
	"database/sql"
	"flashcard_service/pkg/database"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type MySql struct {
	conn string
	db   *sql.DB
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
	return &MySql{
		conn: config.FormatDSN(),
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

func (m *MySql) Query(query string) (*sql.Rows, error) {
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *MySql) Exec(query string, args ...any) (sql.Result, error) {
	r, err := m.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (m *MySql) QueryRow(query string, args ...any) (*sql.Row, error) {
	row := m.db.QueryRow(query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return row, nil
}
