package validate_test

import (
	"errors"
	"testing"

	"github.com/extreme-business/lingo/pkg/validate"
)

func TestStringValidator_Validate(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		v       validate.StringValidator
		args    args
		wantErr bool
	}{
		{
			name:    "should return an error if the string is empty",
			v:       validate.StringValidator{validate.StringMinLength("test", 1)},
			args:    args{s: ""},
			wantErr: true,
		},
		{
			name:    "should return no error if the string is not empty",
			v:       validate.StringValidator{validate.StringMinLength("test", 1)},
			args:    args{s: "a"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("StringValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringNotEmpty(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string is not empty", func(t *testing.T) {
		v := validate.StringNotEmpty("test")
		if err := v("a"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string is empty", func(t *testing.T) {
		v := validate.StringNotEmpty("test")
		if err := v(""); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("error matches field, message and validate.ErrStringRequired", func(t *testing.T) {
		v := validate.StringNotEmpty("test")
		err := v("")

		if !errors.Is(err, validate.ErrEmptyString) {
			t.Errorf("expected error to be %v, got %v", validate.ErrEmptyString, err)
		}

		if err.Field() != "test" {
			t.Errorf("expected field to be test, got %v", err.Field())
		}

		if err.Error() != "test: string is empty" {
			t.Errorf("expected error message to be 'test: string is empty, got %v", err.Error())
		}
	})
}

func TestStringIsUtf8(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string is valid utf8", func(t *testing.T) {
		v := validate.StringUtf8("test")
		if err := v("a"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string is not valid utf8", func(t *testing.T) {
		v := validate.StringUtf8("test")
		if err := v("\xff"); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("error matches field, message and validate.ErrStringIsUtf8", func(t *testing.T) {
		v := validate.StringUtf8("test")
		err := v("\xff")

		if !errors.Is(err, validate.ErrStringIsUtf8) {
			t.Errorf("expected error to be %v, got %v", validate.ErrStringIsUtf8, err)
		}

		if err.Field() != "test" {
			t.Errorf("expected field to be test, got %v", err.Field())
		}

		if err.Error() != "test: string is not valid utf-8" {
			t.Errorf("expected error message to be 'test: string is not valid utf-8', got %v", err.Error())
		}
	})
}

func TestStringMinLength(t *testing.T) {
	t.Run("should return a StringValidatorFunc that return no error if a string is at least n characters long", func(t *testing.T) {
		v := validate.StringMinLength("test", 1)
		if err := v("a"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string is not at least n characters long", func(t *testing.T) {
		v := validate.StringMinLength("test", 1)
		if err := v(""); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("error matches field, message and validate.ErrStringMinLength", func(t *testing.T) {
		v := validate.StringMinLength("test", 1)
		err := v("")

		if !errors.Is(err, validate.ErrStringMinLength) {
			t.Errorf("expected error to be %v, got %v", validate.ErrStringMinLength, err)
		}

		if err.Field() != "test" {
			t.Errorf("expected field to be test, got %v", err.Field())
		}

		if err.Error() != "test: string should be at least 1 characters long" {
			t.Errorf("expected error message to be 'test: string should be at least 1 characters long', got %v", err.Error())
		}
	})
}

func TestStringMaxLength(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string is under n characters", func(t *testing.T) {
		v := validate.StringMaxLength("test", 5)
		if err := v("ab"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string exceeds n characters", func(t *testing.T) {
		v := validate.StringMaxLength("test", 1)
		if err := v("ab"); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("error matches field, message and validate.ErrStringMaxLength", func(t *testing.T) {
		v := validate.StringMaxLength("test", 1)
		err := v("ab")

		if !errors.Is(err, validate.ErrStringMaxLength) {
			t.Errorf("expected error to be %v, got %v", validate.ErrStringMaxLength, err)
		}

		if err.Field() != "test" {
			t.Errorf("expected field to be test, got %v", err.Field())
		}

		if err.Error() != "test: string should be at most 1 characters long" {
			t.Errorf("expected error message to be 'test: string should be at most 1 characters long', got %v", err.Error())
		}
	})
}

func TestStringContainsSpecialChars(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string contains n special characters", func(t *testing.T) {
		v := validate.StringContainsSpecialChars("test", 1)
		if err := v("a!b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		v = validate.StringContainsSpecialChars("test", 2)
		if err := v("a!!b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string does not contain n special characters", func(t *testing.T) {
		v := validate.StringContainsSpecialChars("test", 1)
		if err := v("ab"); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("error matches field, message and validate.ErrStringContainsSpecialChars", func(t *testing.T) {
		v := validate.StringContainsSpecialChars("test", 2)
		err := v("ab!")

		if !errors.Is(err, validate.ErrStringContainsSpecialChars) {
			t.Errorf("expected error to be %v, got %v", validate.ErrStringContainsSpecialChars, err)
		}

		if err.Field() != "test" {
			t.Errorf("expected field to be test, got %v", err.Field())
		}

		if err.Error() != "test: string should contain at least 2 special characters" {
			t.Errorf("expected error message to be 'test: string should contain at least 2 special characters', got %v", err.Error())
		}
	})
}

func TestStringContainsDigits(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string contains n digits", func(t *testing.T) {
		v := validate.StringContainsDigits("test", 1)
		if err := v("a1b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		v = validate.StringContainsDigits("test", 2)
		if err := v("a12b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string does not contain n digits", func(t *testing.T) {
		v := validate.StringContainsDigits("test", 1)
		if err := v("ab"); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("error matches field, message and validate.ErrStringContainsDigits", func(t *testing.T) {
		v := validate.StringContainsDigits("test", 1)
		err := v("ab")

		if !errors.Is(err, validate.ErrStringContainsDigits) {
			t.Errorf("expected error to be %v, got %v", validate.ErrStringContainsDigits, err)
		}

		if err.Field() != "test" {
			t.Errorf("expected field to be test, got %v", err.Field())
		}

		if err.Error() != "test: string should contain at least 1 digits" {
			t.Errorf("expected error message to be 'test: string should contain at least 1 digits', got %v", err.Error())
		}
	})
}

func TestSpecialCharWhitelist(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string contains only special characters in the whitelist", func(t *testing.T) {
		v := validate.SpecialCharWhitelist("test", 1, '!')
		if err := v("a!b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string contains a special character not in the whitelist", func(t *testing.T) {
		v := validate.SpecialCharWhitelist("test", 1, '!')
		if err := v("a@b"); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("error matches field, message and validate.ErrSpecialCharWhitelist", func(t *testing.T) {
		v := validate.SpecialCharWhitelist("test", 1, '!')

		err := v("a@b")

		if !errors.Is(err, validate.ErrSpecialCharWhitelist) {
			t.Errorf("expected error to be %v, got %v", validate.ErrSpecialCharWhitelist, err)
		}

		if err.Field() != "test" {
			t.Errorf("expected field to be test, got %v", err.Field())
		}

		if err.Error() != "test: string contains invalid character '@'" {
			t.Errorf("expected error message to be 'test: string contains invalid character '@', got %v", err.Error())
		}
	})
}
