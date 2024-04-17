package dbtesting

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/dwethmar/lingo/pkg/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func SetupPostgresContainer(ctx context.Context, setup func(connectionString string) error) (*PostgresContainer, error) {
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err := setup(connectionString); err != nil {
		return nil, err
	}

	// 2. Create a snapshot of the database to restore later
	err = container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: container,
		ConnectionString:  connectionString,
	}, nil
}

// SetupTestDB sets up the database connection for testing and
func SetupTestDB(ctx context.Context, t *testing.T, dbc *PostgresContainer) (*sql.DB, func()) {

	db, close, err := database.ConnectPostgres(ctx, dbc.ConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		if err := close(); err != nil {
			t.Fatal(err)
		}
	}
}
