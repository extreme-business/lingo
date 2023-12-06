package validate

import "fmt"

// Error is a validation error.
type Error struct {
	field   string // Field that caused the error.
	Message string // Message with details about the error.
	err     error  // Wrapped error, used for errors.Is.
}

// Field returns the field that caused the error.
func (e *Error) Field() string { return e.field }

// Error returns the error message.
func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.field, e.Message) }

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error { return e.err }
