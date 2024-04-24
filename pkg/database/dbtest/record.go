package dbtest

import (
	"context"
	"database/sql"

	"github.com/dwethmar/lingo/pkg/database"
)

type Query struct {
	Query string
	Args  []interface{}
}

// Recorder records all queries and arguments that are passed to the database.
type Recorder struct {
	db          database.DB
	Queries     []Query
	RowQueries  []Query
	ExecQueries []Query
}

// NewRecorder creates a new query Recorder.
func NewRecorder(db database.DB) *Recorder {
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
