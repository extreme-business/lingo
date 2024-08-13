package validate

import (
	"errors"
	"fmt"
)

// Error is a validation error.
type Error struct {
	field   string // Field that caused the error.
	Message string // Message with details about the error.
	err     error  // Wrapped error, used for errors.Is.
}

func NewError(field, message string, err error) *Error {
	return &Error{field: field, Message: message, err: err}
}

// Field returns the field that caused the error.
func (e *Error) Field() string { return e.field }

// Error returns the error message.
func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.field, e.Message) }

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error { return e.err }

func AssertError(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}

	var e *Error
	if ok := errors.As(err, &e); ok {
		return e, true
	}

	return nil, false
}
