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

// DB is a database connection that uses a DBHandler.
type DB struct {
	handler DBHandler // handler is the database handler.
}

// NewDB creates a new DB with the given sql.DB.
func NewDB(sqlDB *sql.DB, opt ...Option) *DB {
	return NewDBWithHandler(NewSQLDBWrapper(sqlDB), opt...)
}

// NewDBWithHandler creates a new DB with a custom handler.
func NewDBWithHandler(handler DBHandler, opt ...Option) *DB {
	db := &DB{
		handler: handler, // handler is the database handler.
	}

	for _, o := range opt {
		o(db)
	}

	return db
}

// Query executes a query that returns rows, typically a SELECT statement.
func (d *DB) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return d.handler.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (d *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *Row {
	return d.handler.QueryRowContext(ctx, query, args...)
}

// Exec executes a query without returning any rows.
func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.handler.ExecContext(ctx, query, args...)
}

// Begin starts a new transaction.
func (d *DB) Begin(ctx context.Context) (*Tx, error) {
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
