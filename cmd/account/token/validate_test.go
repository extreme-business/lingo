package token_test

import (
	_ "embed"
	"errors"
	"testing"

	"github.com/dwethmar/lingo/cmd/account/token"
)

//go:embed testdata/valid_token.txt
var validToken []byte

func TestNewValidator(t *testing.T) {
	t.Run("should return a new Validator", func(t *testing.T) {
		got := token.NewValidator(nil)
		if got == nil {
			t.Errorf("NewValidator() = %v, want not nil", got)
		}
	})
}

func TestValidator_Validate(t *testing.T) {
	t.Run("valid token should pass", func(t *testing.T) {
		v := token.NewValidator([]byte("secret"))

		c, err := v.Validate(string(validToken))
		if err != nil {
			t.Errorf("Validate() error = %v, want %v", err, nil)
		}

		if c == nil {
			t.Errorf("Validate() = %v, want %v", c, nil)
		}
	})

	t.Run("should return error if token is malformed", func(t *testing.T) {
		v := token.NewValidator([]byte("secret"))

		c, err := v.Validate("invalid")
		if !errors.Is(err, token.ErrTokenMalformed) {
			t.Errorf("Validate() error = %v, want %v", err, token.ErrTokenMalformed)
		}

		if c != nil {
			t.Errorf("Validate() = %v, want %v", c, nil)
		}
	})
}
