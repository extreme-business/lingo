package token

import (
	"errors"
	"testing"
)

const (
	validToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNzEzNDYyMDkwLCJleHAiOjQwODAyMTAwODl9.ke9p7DjmbmB6BSw6-E9a7iDs802b6LVGHUqjqsRMCx8"
)

func TestValidator_Validate(t *testing.T) {
	t.Run("valid token should pass", func(t *testing.T) {
		v := &Validator{
			secretKey: []byte("secret"),
		}

		c, err := v.Validate(validToken)
		if err != nil {
			t.Errorf("Validate() error = %v, want %v", err, nil)
		}

		if c == nil {
			t.Errorf("Validate() = %v, want %v", c, nil)
		}
	})

	t.Run("should return error if token is malformed", func(t *testing.T) {
		v := &Validator{
			secretKey: []byte("secret"),
		}

		c, err := v.Validate("invalid")
		if !errors.Is(err, ErrTokenMalformed) {
			t.Errorf("Validate() error = %v, want %v", err, ErrTokenMalformed)
		}

		if c != nil {
			t.Errorf("Validate() = %v, want %v", c, nil)
		}
	})

}
