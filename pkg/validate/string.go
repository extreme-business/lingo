package validate

import (
	"fmt"
	"unicode"
)

// StringValidatorFunc is a function that validates a string.
type StringValidatorFunc func(string) error

// StringValidator is a list of StringValidatorFunc.
type StringValidator []StringValidatorFunc

func (v StringValidator) Validate(s string) error {
	for _, f := range v {
		if err := f(s); err != nil {
			return err
		}
	}

	return nil
}

// MinLength returns a StringValidatorFunc that checks if a string is at least n characters long.
func MinLength(name string, n int) StringValidatorFunc {
	return func(s string) error {
		if len(s) < n {
			return fmt.Errorf("%s should be at least %d characters long", name, n)
		}
		return nil
	}
}

// MaxLength returns a StringValidatorFunc that checks if a string is at most n characters long.
func MaxLength(name string, n int) StringValidatorFunc {
	return func(s string) error {
		if len(s) > n {
			return fmt.Errorf("%s should be at most %d characters long", name, n)
		}
		return nil
	}
}

// ContainsSpecialChars returns a StringValidatorFunc that checks if a string contains n special characters.
func ContainsSpecialChars(name string, n int) StringValidatorFunc {
	return func(s string) error {
		count := 0
		for _, runeValue := range s {
			if !unicode.IsLetter(runeValue) && !unicode.IsDigit(runeValue) {
				count++
				if count >= n {
					return nil
				}
			}
		}

		return fmt.Errorf("%s does not contain %d special characters", name, n)
	}
}

// ContainsDigits returns a StringValidatorFunc that checks if a string contains n digits.
func ContainsDigits(name string, n int) StringValidatorFunc {
	return func(s string) error {
		count := 0
		for _, runeValue := range s {
			if unicode.IsDigit(runeValue) {
				count++
				if count >= n {
					return nil
				}
			}
		}

		return fmt.Errorf("%s does not contain %d digits", name, n)
	}
}

// SpecialCharWhitelist returns a StringValidatorFunc that checks if a string contains only letters, digits, and exceptions.
func SpecialCharWhitelist(name string, exceptions ...rune) StringValidatorFunc {
	exceptionsMap := make(map[rune]struct{}, len(exceptions))
	for _, r := range exceptions {
		exceptionsMap[r] = struct{}{}
	}

	return func(s string) error {
		for _, runeValue := range s {
			if !unicode.IsLetter(runeValue) && !unicode.IsDigit(runeValue) {
				if _, ok := exceptionsMap[runeValue]; ok {
					continue
				}

				return fmt.Errorf("%s contains an character that is not allowed", name)
			}
		}

		return nil
	}
}
