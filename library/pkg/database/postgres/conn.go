package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // to open postgres with sql.Open
	"github.com/pkg/errors"
)

func NewDBConn(cfg Config) (*DB, error) {
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 5431
	}
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "open sql conn")
	}
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "postgres ping error : (%v)")
	}
	return &DB{db: db}, nil
}

type DB struct {
	db *sql.DB
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}
