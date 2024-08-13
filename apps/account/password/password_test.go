package password_test

import (
	"testing"

	"github.com/extreme-business/lingo/apps/account/password"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash a password", func(t *testing.T) {
		pw := []byte("password")
		hash, err := password.Hash(pw)
		if err != nil {
			t.Fatalf("HashPassword() error = %v", err)
		}

		if err = password.Check(pw, hash); err != nil {
			t.Fatalf("Check() error = %v", err)
		}
	})
}
