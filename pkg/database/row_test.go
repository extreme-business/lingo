package database_test

import (
	"errors"
	"testing"

	"github.com/extreme-business/lingo/pkg/database"
	"github.com/extreme-business/lingo/pkg/database/mock"
)

func TestNewRows(t *testing.T) {
	t.Run("should return a new Rows instance", func(t *testing.T) {
		handler := &mock.RowsHandler{}
		got := database.NewRows(handler)
		if got == nil {
			t.Error("Expected NewRows to return a non-nil Rows instance")
		}
	})
}

func TestRows_Err(t *testing.T) {
	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		anError := errors.New("an error")

		handler := &mock.RowsHandler{
			ErrFunc: func() error {
				return anError
			},
		}

		rows := database.NewRows(handler)
		err := rows.Err()
		if err == nil {
			t.Error("Expected Err to return an error")
		}

		if !errors.Is(err, anError) {
			t.Errorf("Expected Err to return an error, got %v", err)
		}
	})

	t.Run("should return nil if the handler returns nil", func(t *testing.T) {
		handler := &mock.RowsHandler{
			ErrFunc: func() error {
				return nil
			},
		}

		rows := database.NewRows(handler)
		err := rows.Err()
		if err != nil {
			t.Errorf("Expected Err to return nil, got %v", err)
		}
	})
}

func TestRows_Next(t *testing.T) {
	t.Run("should return false if the handler returns false", func(t *testing.T) {
		handler := &mock.RowsHandler{
			NextFunc: func() bool {
				return false
			},
		}

		rows := database.NewRows(handler)
		if rows.Next() {
			t.Error("Expected Next to return false")
		}
	})

	t.Run("should return true if the handler returns true", func(t *testing.T) {
		handler := &mock.RowsHandler{
			NextFunc: func() bool {
				return true
			},
		}

		rows := database.NewRows(handler)
		if !rows.Next() {
			t.Error("Expected Next to return true")
		}
	})
}

func TestRows_Scan(t *testing.T) {
	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		anError := errors.New("an error")

		handler := &mock.RowsHandler{
			ScanFunc: func(_ ...interface{}) error {
				return anError
			},
		}

		rows := database.NewRows(handler)
		err := rows.Scan()
		if err == nil {
			t.Error("Expected Scan to return an error")
		}

		if !errors.Is(err, anError) {
			t.Errorf("Expected Scan to return an error, got %v", err)
		}
	})

	t.Run("should return nil if the handler returns nil", func(t *testing.T) {
		handler := &mock.RowsHandler{
			ScanFunc: func(_ ...interface{}) error {
				return nil
			},
		}

		rows := database.NewRows(handler)
		err := rows.Scan()
		if err != nil {
			t.Errorf("Expected Scan to return nil, got %v", err)
		}
	})
}

func TestRows_Close(t *testing.T) {
	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		anError := errors.New("an error")

		handler := &mock.RowsHandler{
			CloseFunc: func() error {
				return anError
			},
		}

		rows := database.NewRows(handler)
		err := rows.Close()
		if err == nil {
			t.Error("Expected Close to return an error")
		}

		if !errors.Is(err, anError) {
			t.Errorf("Expected Close to return an error, got %v", err)
		}
	})

	t.Run("should return nil if the handler returns nil", func(t *testing.T) {
		handler := &mock.RowsHandler{
			CloseFunc: func() error {
				return nil
			},
		}

		rows := database.NewRows(handler)
		err := rows.Close()
		if err != nil {
			t.Errorf("Expected Close to return nil, got %v", err)
		}
	})
}

func TestNewRow(t *testing.T) {
	t.Run("should return a new Row instance", func(t *testing.T) {
		handler := &mock.RowHandler{}
		got := database.NewRow(handler)
		if got == nil {
			t.Error("Expected NewRow to return a non-nil Row instance")
		}
	})
}

func TestRow_Err(t *testing.T) {
	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		anError := errors.New("an error")

		handler := &mock.RowHandler{
			ErrFunc: func() error {
				return anError
			},
		}

		row := database.NewRow(handler)
		err := row.Err()
		if err == nil {
			t.Error("Expected Err to return an error")
		}

		if !errors.Is(err, anError) {
			t.Errorf("Expected Err to return an error, got %v", err)
		}
	})

	t.Run("should return nil if the handler returns nil", func(t *testing.T) {
		handler := &mock.RowHandler{
			ErrFunc: func() error {
				return nil
			},
		}

		row := database.NewRow(handler)
		err := row.Err()
		if err != nil {
			t.Errorf("Expected Err to return nil, got %v", err)
		}
	})
}

func TestRow_Scan(t *testing.T) {
	t.Run("should return an error if the handler returns an error", func(t *testing.T) {
		anError := errors.New("an error")

		handler := &mock.RowHandler{
			ScanFunc: func(_ ...interface{}) error {
				return anError
			},
		}

		row := database.NewRow(handler)
		err := row.Scan()
		if err == nil {
			t.Error("Expected Scan to return an error")
		}

		if !errors.Is(err, anError) {
			t.Errorf("Expected Scan to return an error, got %v", err)
		}
	})

	t.Run("should return nil if the handler returns nil", func(t *testing.T) {
		handler := &mock.RowHandler{
			ScanFunc: func(_ ...interface{}) error {
				return nil
			},
		}

		row := database.NewRow(handler)
		err := row.Scan()
		if err != nil {
			t.Errorf("Expected Scan to return nil, got %v", err)
		}
	})
}
