package database

import (
	"context"
	"database/sql"
	"fmt"
)

var _ DBHandler = &SQLDBWrapper{}
var _ TXHandler = &SQLTxWrapper{}

// SQLDBWrapper wraps a sql.DB to implement DBHandler.
type SQLDBWrapper struct {
	db *sql.DB
}

// BeginTx implements DBHandler.
func (s *SQLDBWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	sqlTx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return NewTx(sqlTx), nil
}

// ExecContext implements DBHandler.
func (s *SQLDBWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

// QueryContext implements DBHandler.
func (s *SQLDBWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	sqlRows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	if rErr := sqlRows.Err(); rErr != nil {
		return nil, fmt.Errorf("sql rows returned an error: %w", rErr)
	}

	return NewRows(sqlRows), nil
}

// QueryRowContext implements DBHandler.
func (s *SQLDBWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	return NewRow(s.db.QueryRowContext(ctx, query, args...))
}

// NewSQLDBWrapper creates a new SQLDBWrapper.
func NewSQLDBWrapper(db *sql.DB) *SQLDBWrapper {
	return &SQLDBWrapper{db: db}
}

// SQLTxWrapper wraps a sql.Tx to implement TXHandler.
type SQLTxWrapper struct {
	tx *sql.Tx
}

// Commit implements TXHandler.
func (s *SQLTxWrapper) Commit() error {
	return s.tx.Commit()
}

// Rollback implements TXHandler.
func (s *SQLTxWrapper) Rollback() error {
	return s.tx.Rollback()
}

// ExecContext implements TXHandler.
func (s *SQLTxWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.tx.ExecContext(ctx, query, args...)
}

// QueryContext implements TXHandler.
func (s *SQLTxWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	sqlRows, err := s.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	if rErr := sqlRows.Err(); rErr != nil {
		return nil, fmt.Errorf("sql rows returned an error: %w", rErr)
	}

	return NewRows(sqlRows), nil
}

// QueryRowContext implements TXHandler.
func (s *SQLTxWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	return NewRow(s.tx.QueryRowContext(ctx, query, args...))
}

// NewSQLTxWrapper creates a new SQLTxWrapper.
func NewSQLTxWrapper(tx *sql.Tx) *SQLTxWrapper {
	return &SQLTxWrapper{tx: tx}
}
