package mysql

import (
	"context"
	"database/sql"
	"flashcard_service/pkg/database"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

type MySql struct {
	conn         string
	db           *sql.DB
	queryTimeout time.Duration
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
		conn:         config.FormatDSN(),
		queryTimeout: time.Second * time.Duration(queryTimeout),
	}
}

func (m *MySql) Connect() error {
	db, err := sql.Open("mysql", m.conn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to MySQL database")
		return err
	}
	m.db = db
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetConnMaxIdleTime(time.Minute * 3)
	maxOpenCons, _ := strconv.Atoi(os.Getenv("MYSQL_POOL_MAX_OPEN_CONNECTION"))
	maxIdConst, _ := strconv.Atoi(os.Getenv("MYSQL_POOL_MAX_IDLE_CONNECTION"))
	db.SetMaxIdleConns(maxIdConst)
	db.SetMaxOpenConns(maxOpenCons)
	log.Info().Msg("Connected to MySQL database")
	return m.Ping()
}

func (m *MySql) Close() error {
	err := m.db.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close MySQL database connection")
		return err
	}
	log.Info().Msg("Closed MySQL database connection")
	return nil
}

func (m *MySql) Ping() error {
	err := m.db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Failed to ping MySQL database")
		return err
	}
	log.Info().Msg("Pinged MySQL database successfully")
	return nil
}

func (m *MySql) QueryRows(query string, args ...any) (*sql.Rows, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.queryTimeout)
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, cancel, err
	}
	return rows, cancel, nil
}

func (m *MySql) Exec(query string, args ...any) (sql.Result, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.queryTimeout)
	r, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, cancel, err
	}
	return r, cancel, nil
}

func (m *MySql) QueryRow(query string, args ...any) (*sql.Row, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.queryTimeout)
	row := m.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, cancel, row.Err()
	}
	return row, cancel, nil
}
