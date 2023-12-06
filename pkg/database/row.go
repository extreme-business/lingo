package database

import "database/sql"

var _ RowsHandler = &sql.Rows{}
var _ RowHandler = &sql.Row{}

type RowsHandler interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
}

type Rows struct {
	handler RowsHandler
}

func NewRows(handler RowsHandler) *Rows {
	return &Rows{
		handler: handler,
	}
}

// Err returns the error, if any, that was encountered during iteration.
func (r *Rows) Err() error {
	return r.handler.Err()
}

// Next prepares the next result row for reading.
func (r *Rows) Next() bool {
	return r.handler.Next()
}

// Scan copies the columns in the current row into the values pointed at by dest.
func (r *Rows) Scan(dest ...interface{}) error {
	return r.handler.Scan(dest...)
}

// Close closes the Rows, preventing further enumeration.
func (r *Rows) Close() error {
	return r.handler.Close()
}

// RowHandler is a database row and should comply with *sql.Row.
type RowHandler interface {
	Err() error
	Scan(dest ...interface{}) error
}

type Row struct {
	handler RowHandler
}

func NewRow(handler RowHandler) *Row {
	return &Row{
		handler: handler,
	}
}

// Err returns the error, if any, that was encountered during iteration.
func (r *Row) Err() error {
	return r.handler.Err()
}

// Scan copies the columns in the current row into the values pointed at by dest.
func (r *Row) Scan(dest ...interface{}) error {
	return r.handler.Scan(dest...)
}
