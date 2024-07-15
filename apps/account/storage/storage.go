package storage

import "context"

// Repositories is a collection of repositories.
type Repositories struct {
	User         UserRepository
	Organization OrganizationRepository
}

// DBManager is a database manager. It is used to manage the repositories.
type DBManager interface {
	// Op starts a new database operation.
	Op() Repositories
	// BeginTX starts a new transaction, performs the repository operations within that transaction,
	// and commits the transaction if all operations succeed; it rolls back the transaction otherwise.
	BeginOp(ctx context.Context, operation func(context.Context, Repositories) error) error
}

type Pagination struct {
	Limit  int
	Offset int
}

type Direction string

const (
	ASC  Direction = "ASC"
	DESC Direction = "DESC"
)
