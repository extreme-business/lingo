package postgres

import (
	"context"
	"testing"

	"github.com/dwethmar/lingo/cmd/auth/storage/organization"
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/pkg/database"
)

// Run runs the seed.
func Run(t *testing.T, connectionString string, organizations []*organization.Organization, users []*user.User) {
	t.Helper()

	db, err := database.ConnectPostgres(context.Background(), connectionString)
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	for _, o := range organizations {
		if err := Organization(context.Background(), tx, o); err != nil {
			t.Fatal(err)
		}
	}

	for _, u := range users {
		if err := User(context.Background(), tx, u); err != nil {
			t.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
