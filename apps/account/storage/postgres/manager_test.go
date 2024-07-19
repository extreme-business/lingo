package postgres_test

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/extreme-business/lingo/apps/account/storage/postgres"
	"github.com/extreme-business/lingo/pkg/database"
	dbmock "github.com/extreme-business/lingo/pkg/database/mock"
)

// checkIfAllReposAreNotNil checks if all repositories are non-nil.
func checkIfAllReposAreNotNil(t *testing.T, repos storage.Repositories) {
	t.Helper()

	val := reflect.ValueOf(repos)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.IsNil() {
			t.Errorf("Expected New to return a manager with non-nil repositories, but %q is nil", val.Type().Field(i).Name)
		}
	}
}

func TestNew(t *testing.T) {
	t.Run("New repository should have all repositories set", func(t *testing.T) {
		got := postgres.NewManager(database.NewDB(&sql.DB{}))
		if got == nil {
			t.Error("Expected New to return a non-nil manager instance")
			return
		}

		// None of them should be nil
		checkIfAllReposAreNotNil(t, got.Op())
	})

	t.Run("New tx should have all repositories set", func(t *testing.T) {
		db := database.NewDBWithHandler(&dbmock.DBHandler{
			BeginTxFunc: func(_ context.Context, _ *sql.TxOptions) (*database.Tx, error) {
				return database.NewTxWithHandler(&dbmock.TxHandler{
					CommitFunc: func() error { return nil },
				}), nil
			},
		})

		got := postgres.NewManager(db)

		if got == nil {
			t.Error("Expected New to return a non-nil manager instance")
			return
		}

		// None of them should be nil
		err := got.BeginOp(context.TODO(), func(_ context.Context, r storage.Repositories) error {
			checkIfAllReposAreNotNil(t, r)
			return nil
		})

		if err != nil {
			t.Errorf("BeginTX() error = %v, want %v", err, nil)
		}
	})
}
