package validate

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrStringMinLength            = errors.New("string is too short")
	ErrStringMaxLength            = errors.New("string is too long")
	ErrStringContainsSpecialChars = errors.New("string does not contain enough special characters")
	ErrStringContainsDigits       = errors.New("string does not contain enough digits")
	ErrSpecialCharWhitelist       = errors.New("string contains invalid characters")
)

// StringValidatorFunc is a function that validates a string.
type StringValidatorFunc func(string) *Error

// StringValidator is a list of StringValidatorFunc.
type StringValidator []StringValidatorFunc

func (v StringValidator) Validate(s string) *Error {
	for _, f := range v {
		if err := f(s); err != nil {
			return err
		}
	}

	return nil
}

// StringMinLength returns a StringValidatorFunc that checks if a string is at least n characters long.
func StringMinLength(name string, n int) StringValidatorFunc {
	return func(s string) *Error {
		if len(s) < n {
			return &Error{
				field:   name,
				Message: fmt.Sprintf("string should be at least %d characters long", n),
				err:     ErrStringMinLength,
			}
		}
		return nil
	}
}

// StringMaxLength returns a StringValidatorFunc that checks if a string is at most n characters long.
func StringMaxLength(name string, n int) StringValidatorFunc {
	return func(s string) *Error {
		if len(s) > n {
			return &Error{
				field:   name,
				Message: fmt.Sprintf("string should be at most %d characters long", n),
				err:     ErrStringMaxLength,
			}
		}
		return nil
	}
}

// StringContainsSpecialChars returns a StringValidatorFunc that checks if a string contains n special characters.
func StringContainsSpecialChars(name string, n int) StringValidatorFunc {
	return func(s string) *Error {
		count := 0
		for _, runeValue := range s {
			if !unicode.IsLetter(runeValue) && !unicode.IsDigit(runeValue) {
				count++
				if count >= n {
					return nil
				}
			}
		}

		return &Error{
			field:   name,
			Message: fmt.Sprintf("string should contain at least %d special characters", n),
			err:     ErrStringContainsSpecialChars,
		}
	}
}

// StringContainsDigits returns a StringValidatorFunc that checks if a string contains n digits.
func StringContainsDigits(name string, n int) StringValidatorFunc {
	return func(s string) *Error {
		count := 0
		for _, runeValue := range s {
			if unicode.IsDigit(runeValue) {
				count++
				if count >= n {
					return nil
				}
			}
		}

		return &Error{
			field:   name,
			Message: fmt.Sprintf("string should contain at least %d digits", n),
			err:     ErrStringContainsDigits,
		}
	}
}

// SpecialCharWhitelist returns a StringValidatorFunc that checks if a string contains only letters, digits, and exceptions.
func SpecialCharWhitelist(name string, exceptions ...rune) StringValidatorFunc {
	exceptionsMap := make(map[rune]struct{}, len(exceptions))
	for _, r := range exceptions {
		exceptionsMap[r] = struct{}{}
	}

	return func(s string) *Error {
		for _, runeValue := range s {
			if !unicode.IsLetter(runeValue) && !unicode.IsDigit(runeValue) {
				if _, ok := exceptionsMap[runeValue]; ok {
					continue
				}

				return &Error{
					field:   name,
					Message: fmt.Sprintf("string contains invalid character %q", runeValue),
					err:     ErrSpecialCharWhitelist,
				}
			}
		}

		return nil
	}
}
