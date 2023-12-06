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

type Query struct {
	Query string
	Args  []interface{}
}

// Recorder records all queries and arguments that are passed to the database.
type Recorder struct {
	db          DB
	Queries     []Query
	RowQueries  []Query
	ExecQueries []Query
}

func NewRecorder(db DB) *Recorder {
	return &Recorder{db: db}
}

func (d *Recorder) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	d.Queries = append(d.Queries, Query{Query: query, Args: args})
	return d.db.QueryContext(ctx, query, args...)
}

func (d *Recorder) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	d.RowQueries = append(d.RowQueries, Query{Query: query, Args: args})
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *Recorder) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	d.ExecQueries = append(d.ExecQueries, Query{Query: query, Args: args})
	return d.db.ExecContext(ctx, query, args...)
}
