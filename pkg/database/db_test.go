package database_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/dwethmar/lingo/pkg/database"
	"github.com/dwethmar/lingo/pkg/database/mock"
)

func TestNewDB(t *testing.T) {
	t.Run("should return a new DB instance", func(t *testing.T) {
		db := database.NewDB(&sql.DB{})
		if db == nil {
			t.Error("Expected NewDB to return a non-nil DB instance")
		}
	})

	t.Run("should return a new DB instance", func(t *testing.T) {
		handler := &mock.DBHandler{}
		db := database.NewDBWithHandler(handler)

		if db == nil {
			t.Error("Expected NewDB to return a non-nil DB instance")
		}
	})

	t.Run("should return a new DB instance and apply options", func(t *testing.T) {
		var called bool
		handler := &mock.DBHandler{}
		db := database.NewDBWithHandler(handler, func(_ *database.DB) {
			called = true
		})

		if db == nil {
			t.Error("Expected NewDB to return a non-nil DB instance")
		}

		if !called {
			t.Error("Expected option to be called")
		}
	})
}

func TestDB_Query(t *testing.T) {
	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		var anError = errors.New("an error")

		handler := &mock.DBHandler{
			QueryContextFunc: func(_ context.Context, _ string, _ ...interface{}) (*database.Rows, error) {
				return nil, errors.New("an error")
			},
		}

		db := database.NewDBWithHandler(handler)

		rows, err := db.Query(context.Background(), "SELECT * FROM users")
		if err == nil {
			t.Error("Expected Query to return an error")
		}

		if errors.Is(err, anError) {
			t.Errorf("Expected Query to return an error, got %v", err)
		}

		if rows != nil {
			t.Error("Expected Query to return nil")
		}
	})
}

func TestDB_QueryRow(t *testing.T) {
	t.Run("should return row if the handler returns a row", func(t *testing.T) {
		handler := &mock.DBHandler{
			QueryRowContextFunc: func(_ context.Context, _ string, _ ...interface{}) *database.Row {
				return database.NewRow(&database.Row{})
			},
		}

		db := database.NewDBWithHandler(handler)

		row := db.QueryRow(context.Background(), "SELECT * FROM users")
		if row == nil {
			t.Error("Expected QueryRow to return a non-nil Row instance")
		}
	})

	t.Run("should return nil if the handler returns nil", func(t *testing.T) {
		handler := &mock.DBHandler{
			QueryRowContextFunc: func(_ context.Context, _ string, _ ...interface{}) *database.Row {
				return nil
			},
		}

		db := database.NewDBWithHandler(handler)

		row := db.QueryRow(context.Background(), "SELECT * FROM users")
		if row != nil {
			t.Error("Expected QueryRow to return nil")
		}
	})
}

func TestDB_Exec(t *testing.T) {
	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		var anError = errors.New("an error")

		handler := &mock.DBHandler{
			ExecContextFunc: func(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, errors.New("an error")
			},
		}

		db := database.NewDBWithHandler(handler)

		result, err := db.Exec(context.Background(), "UPDATE users SET name = ?", "new_name")
		if err == nil {
			t.Error("Expected Exec to return an error")
		}

		if errors.Is(err, anError) {
			t.Errorf("Expected Exec to return an error, got %v", err)
		}

		if result != nil {
			t.Error("Expected Exec to return nil")
		}
	})

	t.Run("should return a result if the handler returns a result", func(t *testing.T) {
		handler := &mock.DBHandler{
			ExecContextFunc: func(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return &sqlmockResult{}, nil
			},
		}

		db := database.NewDBWithHandler(handler)

		result, err := db.Exec(context.Background(), "UPDATE users SET name = ?", "new_name")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result == nil {
			t.Error("Expected Exec to return a non-nil result")
		}
	})
}

func Test_NewTX(t *testing.T) {
	t.Run("should return a new Tx instance", func(t *testing.T) {
		tx := database.NewTx(&sql.Tx{})

		if tx == nil {
			t.Error("Expected NewTx to return a non-nil Tx instance")
		}
	})

	t.Run("should return a new Tx instance", func(t *testing.T) {
		handler := &mock.TxHandler{}
		tx := database.NewTxWithHandler(handler)

		if tx == nil {
			t.Error("Expected NewTx to return a non-nil Tx instance")
		}
	})
}

func TestTx_Query(t *testing.T) {
	t.Run("should return rows if the handler returns rows", func(t *testing.T) {
		handler := &mock.TxHandler{
			DBHandler: mock.DBHandler{
				QueryContextFunc: func(_ context.Context, _ string, _ ...interface{}) (*database.Rows, error) {
					return &database.Rows{}, nil
				},
			},
		}

		tx := database.NewTxWithHandler(handler)

		rows, err := tx.Query(context.Background(), "SELECT * FROM users")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if rows == nil {
			t.Error("Expected Query to return non-nil rows")
		}
	})

	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		handler := &mock.TxHandler{
			DBHandler: mock.DBHandler{
				QueryContextFunc: func(_ context.Context, _ string, _ ...interface{}) (*database.Rows, error) {
					return nil, errors.New("an error")
				},
			},
		}

		tx := database.NewTxWithHandler(handler)

		rows, err := tx.Query(context.Background(), "SELECT * FROM users")
		if err == nil {
			t.Error("Expected Query to return an error")
		}
		if rows != nil {
			t.Error("Expected Query to return nil")
		}
	})
}

func TestTx_QueryRow(t *testing.T) {
	t.Run("should return row if the handler returns a row", func(t *testing.T) {
		handler := &mock.TxHandler{
			DBHandler: mock.DBHandler{
				QueryRowContextFunc: func(_ context.Context, _ string, _ ...interface{}) *database.Row {
					return &database.Row{}
				},
			},
		}

		tx := database.NewTxWithHandler(handler)

		row := tx.QueryRow(context.Background(), "SELECT * FROM users WHERE id = ?")
		if row == nil {
			t.Error("Expected QueryRow to return a non-nil Row instance")
		}
	})

	t.Run("should return nil if the handler returns nil", func(t *testing.T) {
		handler := &mock.TxHandler{
			DBHandler: mock.DBHandler{
				QueryRowContextFunc: func(_ context.Context, _ string, _ ...interface{}) *database.Row {
					return nil
				},
			},
		}

		tx := database.NewTxWithHandler(handler)

		row := tx.QueryRow(context.Background(), "SELECT * FROM users WHERE id = ?")
		if row != nil {
			t.Error("Expected QueryRow to return nil")
		}
	})
}

func TestTx_Exec(t *testing.T) {
	t.Run("should return result if the handler returns a result", func(t *testing.T) {
		handler := &mock.TxHandler{
			DBHandler: mock.DBHandler{
				ExecContextFunc: func(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
					return &sqlmockResult{}, nil
				},
			},
		}

		tx := database.NewTxWithHandler(handler)

		result, err := tx.Exec(context.Background(), "UPDATE users SET name = ?", "new_name")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result == nil {
			t.Error("Expected Exec to return a non-nil result")
		}
	})

	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		handler := &mock.TxHandler{
			DBHandler: mock.DBHandler{
				ExecContextFunc: func(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
					return nil, errors.New("an error")
				},
			},
		}

		tx := database.NewTxWithHandler(handler)

		result, err := tx.Exec(context.Background(), "UPDATE users SET name = ?", "new_name")
		if err == nil {
			t.Error("Expected Exec to return an error")
		}
		if result != nil {
			t.Error("Expected Exec to return nil")
		}
	})
}

func TestTx_Commit(t *testing.T) {
	t.Run("should commit the transaction successfully", func(t *testing.T) {
		handler := &mock.TxHandler{
			CommitFunc: func() error {
				return nil
			},
		}

		tx := database.NewTxWithHandler(handler)

		err := tx.Commit()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		handler := &mock.TxHandler{
			CommitFunc: func() error {
				return errors.New("an error")
			},
		}

		tx := database.NewTxWithHandler(handler)

		err := tx.Commit()
		if err == nil {
			t.Error("Expected Commit to return an error")
		}
	})
}

func TestTx_Rollback(t *testing.T) {
	t.Run("should rollback the transaction successfully", func(t *testing.T) {
		handler := &mock.TxHandler{
			RollbackFunc: func() error {
				return nil
			},
		}

		tx := database.NewTxWithHandler(handler)

		err := tx.Rollback()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		handler := &mock.TxHandler{
			RollbackFunc: func() error {
				return errors.New("an error")
			},
		}

		tx := database.NewTxWithHandler(handler)

		err := tx.Rollback()
		if err == nil {
			t.Error("Expected Rollback to return an error")
		}
	})
}

// sqlmockResult is a mock implementation of sql.Result for testing purposes.
type sqlmockResult struct{}

func (r *sqlmockResult) LastInsertId() (int64, error) { return 0, nil }
func (r *sqlmockResult) RowsAffected() (int64, error) { return 1, nil }
