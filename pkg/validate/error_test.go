package validate_test

import (
	"errors"
	"testing"

	"github.com/extreme-business/lingo/pkg/validate"
)

func TestError_Field(t *testing.T) {
	t.Run("should return the field that caused the error", func(t *testing.T) {
		e := validate.NewError("test", "message", nil)
		if got := e.Field(); got != "test" {
			t.Errorf("Error.Field() = %v, want %v", got, "test")
		}
	})
}

func TestError_Error(t *testing.T) {
	t.Run("should return the error message", func(t *testing.T) {
		e := validate.NewError("test", "message", nil)
		if got := e.Error(); got != "test: message" {
			t.Errorf("Error.Error() = %v, want %v", got, "test: message")
		}
	})
}

func TestError_Unwrap(t *testing.T) {
	t.Run("should return the wrapped error", func(t *testing.T) {
		err := errors.New("wrapped")
		e := validate.NewError("test", "message", err)
		got := e.Unwrap()

		if !errors.Is(got, err) {
			t.Errorf("Error.Unwrap() = %v, want %v", got, err)
		}
	})
}

func TestAssertError(t *testing.T) {
	t.Run("should return the error and true if it is a *Error", func(t *testing.T) {
		err := validate.NewError("test", "message", errors.New("wrapped"))
		got, ok := validate.AssertError(err)

		if !ok {
			t.Errorf("AssertError() = %v, want %v", ok, true)
		}

		if !errors.Is(err, got) {
			t.Errorf("AssertError() = %v, want %v", got, err)
		}
	})

	t.Run("should return nil and false if it is not a *Error", func(t *testing.T) {
		err := errors.New("error")
		got, ok := validate.AssertError(err)

		if ok {
			t.Errorf("AssertError() = %v, want %v", ok, false)
		}

		if got != nil {
			t.Errorf("AssertError() = %v, want %v", got, nil)
		}
	})
}
