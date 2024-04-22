// database is a package that contains the database interface. It is used to make it easier to use transactions.
// TX is a transaction interface and complies with the DB interface.
package database

import (
	"context"
	"database/sql"
)

var (
	_ DB = (*sql.DB)(nil)
	_ DB = (*sql.Tx)(nil)
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
