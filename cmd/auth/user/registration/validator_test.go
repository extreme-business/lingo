package registration

import (
	"testing"
)

func TestNewRegistrationValidator(t *testing.T) {
	t.Run("should return a new validator", func(t *testing.T) {
		v := NewRegistrationValidator()
		if v == nil {
			t.Error("expected a new validator")
			return
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
