package database

import (
	"context"
	"database/sql"
	"fmt"
)

// SetupDatabase sets up the database connection.
func Connect(ctx context.Context, dataSourceName string) (DB, func() error, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open db: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("could not ping db: %w", err)
	}

	return db, db.Close, nil
}
