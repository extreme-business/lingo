package validate

type Error struct {
	field string
	err   error
}

// Field returns the field that caused the error.
func (e *Error) Field() string { return e.field }

// Error returns the error message.
func (e *Error) Error() string { return e.err.Error() }

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error { return e.err }
