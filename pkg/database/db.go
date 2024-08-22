// database is a package that contains the database interface. It is used to make it easier to use transactions.
// TX is a transaction interface and complies with the DB interface.
package database

import (
	"context"
	"database/sql"
	"errors"
)

var (
	// ErrTransactorNotSet indicates that the transactor is not set.
	ErrTransactorNotSet = errors.New("transactor not set")
)

// DBHandler is a database interface and should comply with *sql.DB.
type DBHandler interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error)
}

// DBWrapper is wrapper for a database handler and should comply with *sql.DB.
type DBWrapper struct {
	handler DBHandler // handler is the database handler.
}

// NewDBWrapper creates a new DB with the given sql.DB.
func NewDBWrapper(db *sql.DB) *DBWrapper {
	return NewDBWithHandler(NewSQLDBWrapper(db))
}

// NewDBWithHandler creates a new DB with a custom handler.
func NewDBWithHandler(handler DBHandler) *DBWrapper {
	db := &DBWrapper{
		handler: handler, // handler is the database handler.
	}

	return db
}

// Query executes a query that returns rows, typically a SELECT statement.
func (d *DBWrapper) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return d.handler.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (d *DBWrapper) QueryRow(ctx context.Context, query string, args ...interface{}) *Row {
	return d.handler.QueryRowContext(ctx, query, args...)
}

// Exec executes a query without returning any rows.
func (d *DBWrapper) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.handler.ExecContext(ctx, query, args...)
}

// Begin starts a new transaction.
func (d *DBWrapper) Begin(ctx context.Context) (*Tx, error) {
	return d.handler.BeginTx(ctx, nil)
}

// TXHandler is a database transaction and should comply with *sql.Tx.
type TXHandler interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
}

// Tx is a transaction that uses a TXHandler.
type Tx struct {
	handler TXHandler
}

// NewTx creates a new Tx with the given sql.Tx.
func NewTx(sqlTx *sql.Tx) *Tx {
	return NewTxWithHandler(NewSQLTxWrapper(sqlTx))
}

// NewTxWithHandler creates a new Tx with a custom handler.
func NewTxWithHandler(handler TXHandler) *Tx {
	return &Tx{
		handler: handler,
	}
}

// Query executes a query that returns rows, typically a SELECT statement.
func (t *Tx) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return t.handler.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (t *Tx) QueryRow(ctx context.Context, query string, args ...interface{}) *Row {
	return t.handler.QueryRowContext(ctx, query, args...)
}

// Exec executes a query without returning any rows.
func (t *Tx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.handler.ExecContext(ctx, query, args...)
}

// Commit commits the transaction.
func (t *Tx) Commit() error {
	return t.handler.Commit()
}

// Rollback rolls back the transaction.
func (t *Tx) Rollback() error {
	return t.handler.Rollback()
}
