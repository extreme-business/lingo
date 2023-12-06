package database

import (
	"context"
	"log"
)

// Factory is a function that initializes a repository with a database connection.
type Factory[T any] func(c Conn) T

// Manager manages database connections and transactions for a specific type of repository.
type Manager[T any] struct {
	// db is the database connection.
	db *DB
	// factory is a function that initializes an operation with a database connection.
	factory Factory[T]
	// failingRollbackHandler is a function that is called when a transaction failed to roll back.
	failingRollbackHandler func(ctx context.Context, err error)
}

// NewManager creates a new Manager instance, initializing it with a database connection,
// a transaction manager, and a slice of repository registration functions.
func NewManager[T any](db *DB, register Factory[T]) *Manager[T] {
	return &Manager[T]{
		db:      db,
		factory: register,
		failingRollbackHandler: func(_ context.Context, err error) {
			log.Printf("failed to rollback transaction: %v", err)
		},
	}
}

// SetFailingRollbackHandler sets the handler function that is called when a transaction fails to roll back.
func (m *Manager[T]) SetFailingRollbackHandler(handler func(ctx context.Context, err error)) {
	m.failingRollbackHandler = handler
}

// Op initializes a new operation with the database connection.
func (m *Manager[T]) Op() T {
	return m.factory(m.db)
}

// BeginOp starts a new operation with a transaction and commits the transaction if all operations succeed; it rolls back the transaction otherwise.
func (m *Manager[T]) BeginOp(ctx context.Context, operation func(ctx context.Context, r T) error) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Deferred function to handle rolling back the transaction in case of an error.
	// Because the state of the transaction is uncertain after a commit failure, it is best practice to explicitly call Rollback().
	defer func() {
		// Check if the transaction is being closed with an unhandled error.
		if p := recover(); p != nil {
			if rErr := tx.Rollback(); rErr != nil {
				m.failingRollbackHandler(ctx, rErr)
			}
			panic(p) // re-throw panic after rollback
		}

		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				m.failingRollbackHandler(ctx, rErr)
			}
		}
	}()

	// Initialize the repositories with the transaction.
	r := m.factory(tx)

	// Perform the passed operation; if it fails, return the error to trigger a rollback.
	if err = operation(ctx, r); err != nil {
		return err
	}

	// Commit the transaction if all operations were successful.
	err = tx.Commit() // set the err variable to the result of the commit operation. If it fails, the deferred function will handle the rollback.
	if err != nil {
		return err
	}

	return nil
}
