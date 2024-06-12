package mock

import (
	"context"
	"database/sql"

	"github.com/extreme-business/lingo/pkg/database"
)

// DBHandler is a mock implementation of the database.DBHandler interface.
type DBHandler struct {
	QueryContextFunc    func(ctx context.Context, query string, args ...interface{}) (*database.Rows, error)
	QueryRowContextFunc func(ctx context.Context, query string, args ...interface{}) *database.Row
	ExecContextFunc     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	BeginTxFunc         func(ctx context.Context, opts *sql.TxOptions) (*database.Tx, error)
}

func (m *DBHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (*database.Rows, error) {
	if m.QueryContextFunc == nil {
		panic("call to unimplemented method QueryContext on Handler")
	}
	return m.QueryContextFunc(ctx, query, args...)
}

func (m *DBHandler) QueryRowContext(ctx context.Context, query string, args ...interface{}) *database.Row {
	if m.QueryRowContextFunc == nil {
		panic("call to unimplemented method QueryRowContext on Handler")
	}
	return m.QueryRowContextFunc(ctx, query, args...)
}

func (m *DBHandler) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if m.ExecContextFunc == nil {
		panic("call to unimplemented method ExecContext on Handler")
	}
	return m.ExecContextFunc(ctx, query, args...)
}

func (m *DBHandler) BeginTx(ctx context.Context, opts *sql.TxOptions) (*database.Tx, error) {
	if m.BeginTxFunc == nil {
		panic("call to unimplemented method BeginTx on Handler")
	}
	return m.BeginTxFunc(ctx, opts)
}

// TxHandler is a mock implementation of the database.TxHandler interface.
type TxHandler struct {
	DBHandler
	CommitFunc   func() error
	RollbackFunc func() error
}

func (m *TxHandler) Commit() error {
	if m.CommitFunc == nil {
		panic("call to unimplemented method Commit on TxHandler")
	}
	return m.CommitFunc()
}

func (m *TxHandler) Rollback() error {
	if m.RollbackFunc == nil {
		panic("call to unimplemented method Rollback on TxHandler")
	}
	return m.RollbackFunc()
}

// Tx is a mock implementation of the database.Tx interface.
type Transactor struct {
	BeginTxFunc func(ctx context.Context, opts *sql.TxOptions) (*database.Tx, error)
}

func (m *Transactor) BeginTx(ctx context.Context, opts *sql.TxOptions) (*database.Tx, error) {
	if m.BeginTxFunc == nil {
		panic("call to unimplemented method BeginTx on Transactor")
	}
	return m.BeginTxFunc(ctx, opts)
}
