package domain

import (
	"context"

	userStorage "github.com/extreme-business/lingo/apps/account/storage/user"
	organizationStorage "github.com/extreme-business/lingo/apps/organization/storage/organization"
	"github.com/extreme-business/lingo/pkg/database"
)

type Repositories struct {
	UserRepo         *userStorage.UserStorage
	OrganizationRepo *organizationStorage.OrganizationStorage
}

// Repositories is a collection of repositories.
type Domain struct {
	UserReader *UserReader
	UserWriter *UserWriter
	// OrganizationReader OrganizationReader
	// OrganizationWriter OrganizationWriter
}

// Manager is a database manager. It is used to manage the repositories.
type Manager interface {
	// Op starts a new database operation.
	Op() Domain
	// BeginOp starts a new transaction, performs the repository operations within that transaction,
	// and commits the transaction if all operations succeed; it rolls back the transaction otherwise.
	BeginOp(ctx context.Context, operation func(context.Context, Domain) error) error
}

// NewManager creates a new manager for storage.
func NewManager(db *database.DB, f func(c database.Conn) *Repositories) *database.Manager[Domain] {
	return database.NewManager(db, func(c database.Conn) Domain {
		repos := f(c)
		return Domain{
			UserReader: NewUserReader(repos.UserRepo),
			UserWriter: NewUserWriter(repos.UserRepo),
		}
	})
}
