package registration

import (
	"reflect"
	"testing"

	"github.com/dwethmar/lingo/pkg/validate"
)

func TestNewRegistrationValidator(t *testing.T) {
	t.Run("should return a new validator", func(t *testing.T) {
		v := NewRegistrationValidator()
		if v == nil {
			t.Error("expected a new validator")
			return
		}

		expected := &RegistrationValidator{
			usernameValidator: validate.StringValidator{
				validate.MaxLength("username", 50),
				validate.MinLength("username", 3),
				validate.SpecialCharWhitelist("username", '_', '-'),
			},
			emailValidator: validate.StringValidator{
				validate.MaxLength("email", 50),
				validate.MinLength("email", 3),
			},
			passwordValidator: validate.StringValidator{
				validate.MinLength("password", 8),
				validate.MaxLength("password", 50),
				validate.ContainsSpecialChars("password", 1),
				validate.ContainsDigits("password", 1),
			},
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("NewRegistrationValidator() = %v, want %v", v, expected)
		}
	})
}

func TestRegistrationValidator_Validate(t *testing.T) {
	t.Run("should return no error if the registration is valid", func(t *testing.T) {
		v := NewRegistrationValidator()
		err := v.Validate(Registration{
			Username: "test-username",
			Email:    "test@test.com",
			Password: "test-password1",
		})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return an error if the registration is invalid", func(t *testing.T) {
		v := NewRegistrationValidator()
		err := v.Validate(Registration{
			Username: "a", // too short
			Email:    "test@test.com",
			Password: "test-password1",
		})

		if err == nil {
			t.Error("expected an error")
		}
	})
}
