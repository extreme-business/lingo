// database is a package that contains the database interface. It is used to make it easier to use transactions.
// TX is a transaction interface and complies with the DB interface.
package database

import (
	"database/sql"
	"testing"
)

func TestNew(t *testing.T) {
	db := &sql.DB{}

	want := &Transactor{
		db: db,
	}

	if got := New(db); got.db != want.db {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func TestTransactor_DB(t *testing.T) {
	db := &sql.DB{}

	want := &Transactor{
		db: db,
	}

	if got := New(db); got.DB() != want.db {
		t.Errorf("DB() = %v, want %v", got, want)
	}
}
