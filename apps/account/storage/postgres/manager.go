package postgres

import (
	"database/sql"

	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/extreme-business/lingo/apps/account/storage/postgres/organization"
	"github.com/extreme-business/lingo/apps/account/storage/postgres/user"
	"github.com/extreme-business/lingo/pkg/database"
)

func Factory(c database.Conn) storage.Repositories {
	return storage.Repositories{
		User:         user.New(c),
		Organization: organization.New(c),
	}
}

// NewManager creates a new manager for storage.
func NewManager(db *sql.DB) *database.Manager[storage.Repositories] {
	return database.NewManager(database.NewDBWrapper(db), Factory)
}

func NewManagerWithWrapper(db *database.DBWrapper) *database.Manager[storage.Repositories] {
	return database.NewManager(db, Factory)
}
