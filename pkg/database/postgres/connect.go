package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // load the postgres driver
)

// SetupDatabase sets up the database connection.
func Connect(ctx context.Context, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
