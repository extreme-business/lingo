package dbtest

import (
	"context"
	"database/sql"

	"github.com/dwethmar/lingo/pkg/database"
)

var _ database.Conn = &Recorder{}

type Query struct {
	Query string
	Args  []interface{}
}

// Recorder records all queries and arguments that are passed to the database.
type Recorder struct {
	dbConn      database.Conn
	Queries     []Query
	RowQueries  []Query
	ExecQueries []Query
}

// NewRecorder creates a new query Recorder.
func NewRecorder(dbConn database.Conn) *Recorder {
	return &Recorder{dbConn: dbConn}
}

func (d *Recorder) Query(ctx context.Context, query string, args ...interface{}) (*database.Rows, error) {
	d.Queries = append(d.Queries, Query{Query: query, Args: args})
	return d.dbConn.Query(ctx, query, args...)
}

func (d *Recorder) QueryRow(ctx context.Context, query string, args ...interface{}) *database.Row {
	d.RowQueries = append(d.RowQueries, Query{Query: query, Args: args})
	return d.dbConn.QueryRow(ctx, query, args...)
}

func (d *Recorder) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	d.ExecQueries = append(d.ExecQueries, Query{Query: query, Args: args})
	return d.dbConn.Exec(ctx, query, args...)
}
