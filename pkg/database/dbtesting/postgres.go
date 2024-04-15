package dbtesting

import (
	"context"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

type PostgresContainer struct {
	URL               string
	postgresContainer *postgres.PostgresContainer
}

func (p *PostgresContainer) Terminate(ctx context.Context) error {
	return p.postgresContainer.Terminate(ctx)
}

func (p *PostgresContainer) Restore(ctx context.Context) error {
	return p.postgresContainer.Restore(ctx, postgres.WithSnapshotName("initial"))
}

func SetupPostgres(ctx context.Context, setup func(dbURL string) error) (*PostgresContainer, error) {
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.RunContainer(
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

	dbURL, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err := setup(dbURL); err != nil {
		return nil, err
	}

	// Create a snapshot of the database to restore later
	err = postgresContainer.Snapshot(ctx, postgres.WithSnapshotName("initial"))
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		URL:               dbURL,
		postgresContainer: postgresContainer,
	}, nil
}
