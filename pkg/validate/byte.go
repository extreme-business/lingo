package validate

import (
	"errors"
	"fmt"
)

var (
	ErrBytesMinLength = errors.New("bytes is too short")
	ErrBytesMaxLength = errors.New("bytes is too long")
)

// StringValidatorFunc is a function that validates a string.
type BytesValidatorFunc func([]byte) *Error

// StringValidator is a list of StringValidatorFunc.
type BytesValidator []BytesValidatorFunc

func (v BytesValidator) Validate(s []byte) *Error {
	for _, f := range v {
		if err := f(s); err != nil {
			return err
		}
	}

	return nil
}

// ByteMinLength returns a StringValidatorFunc that checks if a string is at least n characters long.
func BytesMinLength(field string, l int) BytesValidatorFunc {
	return func(s []byte) *Error {
		if len(s) < l {
			return &Error{
				field:   field,
				Message: fmt.Sprintf("bytes should be at least %d characters long", l),
				err:     ErrBytesMinLength,
			}
		}
		return nil
	}
}

// StringMaxLength returns a StringValidatorFunc that checks if a string is at most n characters long.
func ByteMaxLength(field string, l int) BytesValidatorFunc {
	return func(s []byte) *Error {
		if len(s) > l {
			return &Error{
				field:   field,
				Message: fmt.Sprintf("bytes should be at most %d characters long", l),
				err:     ErrBytesMaxLength,
			}
		}
		return nil
	}
}
