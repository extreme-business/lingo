package validate

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUUIDIsNil = errors.New("uuid is nil")
)

// UUIDValidatorFunc is a function that validates a UUID.
type UUIDValidatorFunc func(uuid.UUID) *Error

// UUIDValidator is a list of UUIDValidatorFunc.
type UUIDValidator []UUIDValidatorFunc

func (v UUIDValidator) Validate(s uuid.UUID) *Error {
	for _, f := range v {
		if err := f(s); err != nil {
			return err
		}
	}

	return nil
}

// UUIDIsNotNil validates that a UUID is not nil.
func UUIDIsNotNil(field string) UUIDValidatorFunc {
	return func(s uuid.UUID) *Error {
		if s == uuid.Nil {
			return &Error{
				field:   field,
				Message: "should not be nil",
				err:     ErrUUIDIsNil,
			}
		}
		return nil
	}
}
