// database is a package that contains the database interface. It is used to make it easier to use transactions.
// TX is a transaction interface and complies with the DB interface.
package database

import (
	"context"
	"database/sql"
	"fmt"
)

var (
	_ DB = (*sql.DB)(nil)
	_ DB = (*sql.Tx)(nil)
	_ TX = (*sql.Tx)(nil)
)

// DB is a database interface.
//
// It is used to make it easier to mock the database and
// to make it easier to use transactions.
type DB interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// TX is a transaction interface.
type TX interface {
	DB
	Commit() error
	Rollback() error
}

// Transactor is a transaction helper that can be used to start transactions.
type Transactor struct {
	db *sql.DB
}

// New creates a new Transactor.
func New(db *sql.DB) *Transactor {
	return &Transactor{
		db: db,
	}
}

// DB returns the database.
func (t *Transactor) DB() DB {
	return t.db
}

// Begin starts a transaction.
func (t *Transactor) Begin(ctx context.Context, opts *sql.TxOptions) (TX, error) {
	tx, err := t.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}

	return tx, nil
}
