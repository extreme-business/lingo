package dbtest

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	dbpostgres "github.com/extreme-business/lingo/pkg/database/postgres"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	occurrenceToWaitFor = 2
	startupTimeout      = 5 * time.Second
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func setupPostgresContainer(ctx context.Context, dbName string, setup func(connectionString string) error) (*PostgresContainer, error) {
	dbUser := "postgres"
	dbPassword := "postgres"

	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(occurrenceToWaitFor).
				WithStartupTimeout(startupTimeout)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err = setup(connectionString); err != nil {
		return nil, err
	}

	// Create a snapshot of the database to restore later
	err = container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: container,
		ConnectionString:  connectionString,
	}, nil
}

// SetupPostgres sets up a test database and runs the provided setup function.
// It also sets up a cleanup function to terminate the container after the test is complete.
func SetupPostgres(t *testing.T, dbName string, setup func(connectionString string) error) *PostgresContainer {
	t.Helper()

	dbc, dbErr := setupPostgresContainer(context.Background(), dbName, setup)

	if dbErr != nil {
		t.Fatal(dbErr)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := dbc.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	return dbc
}

// Connect connects to the test database and returns the connection.
// It also sets up a cleanup function to close the connection after the test is complete.
func Connect(ctx context.Context, t *testing.T, connectionString string) *sql.DB {
	t.Helper()

	db, err := dbpostgres.Connect(ctx, connectionString)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err = db.Close(); err != nil {
			t.Fatalf("failed to close db: %s", err)
		}
	})

	return db
}
