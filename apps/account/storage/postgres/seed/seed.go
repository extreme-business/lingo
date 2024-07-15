package seed

import (
	"context"
	"testing"

	"github.com/extreme-business/lingo/cmd/account/storage"
	"github.com/extreme-business/lingo/pkg/database/postgres"
)

type State struct {
	Organizations []*storage.Organization
	Users         []*storage.User
}

// Run runs the seed.
func Run(t *testing.T, connectionString string, s State) {
	t.Helper()

	db, err := postgres.Connect(context.Background(), connectionString)
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	for _, o := range s.Organizations {
		if err = InsertOrganization(context.Background(), tx, o); err != nil {
			t.Fatal(err)
		}
	}

	for _, u := range s.Users {
		if err = InsertUser(context.Background(), tx, u); err != nil {
			t.Fatal(err)
		}
	}

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
