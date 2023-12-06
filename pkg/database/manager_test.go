package database_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/dwethmar/lingo/pkg/database"
	"github.com/dwethmar/lingo/pkg/database/mock"
)

// TestNewManager tests the NewManager function for proper initialization.
func TestNewManager(t *testing.T) {
	t.Run("should return a new Manager instance", func(t *testing.T) {
		db := database.NewDBWithHandler(&mock.DBHandler{})

		manager := database.NewManager[interface{}](db, nil)

		if manager == nil {
			t.Error("Expected NewManager to return a non-nil manager instance")
		}
	})
}

func TestSetFailingRollbackHandler(t *testing.T) {
	t.Run("should call the failing rollback handler if the transaction rollback fails", func(t *testing.T) {
		calledRollbackHandler := false

		db := database.NewDBWithHandler(&mock.DBHandler{
			BeginTxFunc: func(_ context.Context, _ *sql.TxOptions) (*database.Tx, error) {
				return database.NewTxWithHandler(&mock.TxHandler{
					RollbackFunc: func() error {
						return errors.New("rollback failed")
					},
				}), nil
			},
		})

		manager := database.NewManager(db, func(_ database.Conn) int { return 0 })
		manager.SetFailingRollbackHandler(func(_ context.Context, _ error) {
			calledRollbackHandler = true
		})

		if err := manager.BeginOp(context.Background(), func(_ context.Context, _ int) error {
			return sql.ErrNoRows
		}); err == nil {
			t.Error("Expected BeginTX to return an error")
		}

		if !calledRollbackHandler {
			t.Error("Expected BeginTX to call the failing rollback handler")
		}
	})
}

// TestRepos tests the Repos method to ensure it initializes repositories correctly.
func TestOp(t *testing.T) {
	t.Run("should initialize repositories", func(t *testing.T) {
		calledRegister := false

		register := func(_ database.Conn) int {
			calledRegister = true
			return 0
		}

		manager := database.NewManager(nil, register)

		_ = manager.Op()

		if !calledRegister {
			t.Error("Expected Repos to call the register function")
		}
	})
}

// TestBeginOp tests the BeginOp method to ensure it starts a new operation with a transaction.
func TestBeginOp(t *testing.T) {
	t.Run("should return an error if the transaction fails to begin", func(t *testing.T) {
		beginErr := errors.New("begin failed")

		db := database.NewDBWithHandler(&mock.DBHandler{
			BeginTxFunc: func(_ context.Context, _ *sql.TxOptions) (*database.Tx, error) {
				return nil, beginErr
			},
		})

		register := func(_ database.Conn) int { return 0 }
		manager := database.NewManager(db, register)

		err := manager.BeginOp(context.Background(), func(_ context.Context, _ int) error {
			return nil
		})

		if err == nil {
			t.Error("Expected BeginTX to return an error")
		}

		if !errors.Is(err, beginErr) {
			t.Errorf("Expected BeginTX to return an error, got %v", err)
		}
	})

	t.Run("should rollback the transaction if the operation return error", func(t *testing.T) {
		calledRollback := false

		db := database.NewDBWithHandler(&mock.DBHandler{
			BeginTxFunc: func(_ context.Context, _ *sql.TxOptions) (*database.Tx, error) {
				return database.NewTxWithHandler(&mock.TxHandler{
					RollbackFunc: func() error {
						calledRollback = true
						return nil
					},
				}), nil
			},
		})

		register := func(_ database.Conn) int { return 0 }
		manager := database.NewManager(db, register)

		err := manager.BeginOp(context.Background(), func(_ context.Context, _ int) error {
			return sql.ErrNoRows
		})

		if err == nil {
			t.Error("Expected BeginTX to return an error")
		}

		if !calledRollback {
			t.Error("Expected BeginTX to call Rollback")
		}
	})

	t.Run("if operations panics and rollback fails it should call failingRollbackHandler", func(t *testing.T) {
		calledSetFailingRollbackHandler := false

		db := database.NewDBWithHandler(&mock.DBHandler{
			BeginTxFunc: func(_ context.Context, _ *sql.TxOptions) (*database.Tx, error) {
				return database.NewTxWithHandler(&mock.TxHandler{
					RollbackFunc: func() error {
						return errors.New("rollback failed")
					},
				}), nil
			},
		})

		register := func(_ database.Conn) int { return 0 }
		manager := database.NewManager(db, register)
		manager.SetFailingRollbackHandler(func(_ context.Context, _ error) {
			calledSetFailingRollbackHandler = true
		})

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected BeginTX to panic")
			}

			if !calledSetFailingRollbackHandler {
				t.Error("Expected BeginTX to call the failing rollback handler")
			}
		}()

		err := manager.BeginOp(context.Background(), func(_ context.Context, _ int) error {
			panic("panic")
		})

		if err == nil {
			t.Error("Expected BeginTX to return an error")
		}
	})

	t.Run("if operations return and error and rollback fails it should call failingRollbackHandler", func(t *testing.T) {
		calledSetFailingRollbackHandler := false

		db := database.NewDBWithHandler(&mock.DBHandler{
			BeginTxFunc: func(_ context.Context, _ *sql.TxOptions) (*database.Tx, error) {
				return database.NewTxWithHandler(&mock.TxHandler{
					RollbackFunc: func() error {
						return errors.New("rollback failed")
					},
				}), nil
			},
		})

		register := func(_ database.Conn) int { return 0 }
		manager := database.NewManager(db, register)
		manager.SetFailingRollbackHandler(func(_ context.Context, _ error) {
			calledSetFailingRollbackHandler = true
		})

		err := manager.BeginOp(context.Background(), func(_ context.Context, _ int) error {
			return errors.New("error")
		})

		if err == nil {
			t.Error("Expected BeginTX to return an error")
		}

		if !calledSetFailingRollbackHandler {
			t.Error("Expected BeginTX to call the failing rollback handler")
		}
	})

	t.Run("should rollback the transaction if the commit fails", func(t *testing.T) {
		calledRollbackHandler := false

		db := database.NewDBWithHandler(&mock.DBHandler{
			BeginTxFunc: func(_ context.Context, _ *sql.TxOptions) (*database.Tx, error) {
				return database.NewTxWithHandler(&mock.TxHandler{
					CommitFunc: func() error {
						return errors.New("commit failed")
					},
					RollbackFunc: func() error {
						calledRollbackHandler = true
						return nil
					},
				}), nil
			},
		})

		register := func(_ database.Conn) int { return 0 }

		manager := database.NewManager(db, register)

		err := manager.BeginOp(context.Background(), func(_ context.Context, _ int) error {
			return nil
		})

		if err == nil {
			t.Error("Expected BeginTX to return an error")
		}

		if !calledRollbackHandler {
			t.Error("Expected BeginTX to call the failing rollback handler")
		}
	})

	t.Run("should commit the transaction if the operation succeeds", func(t *testing.T) {
		calledCommit := false

		db := database.NewDBWithHandler(&mock.DBHandler{
			BeginTxFunc: func(_ context.Context, _ *sql.TxOptions) (*database.Tx, error) {
				return database.NewTxWithHandler(&mock.TxHandler{
					CommitFunc: func() error {
						calledCommit = true
						return nil
					},
				}), nil
			},
		})

		register := func(_ database.Conn) int { return 0 }
		manager := database.NewManager(db, register)

		err := manager.BeginOp(context.Background(), func(_ context.Context, _ int) error {
			return nil
		})

		if err != nil {
			t.Errorf("Expected BeginTX to return nil, got %v", err)
		}

		if !calledCommit {
			t.Error("Expected BeginTX to call Commit")
		}
	})
}
