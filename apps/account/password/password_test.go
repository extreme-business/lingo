package password_test

import (
	"testing"

	"github.com/extreme-business/lingo/apps/account/password"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash a password", func(t *testing.T) {
		pw := "password"
		hash, err := password.Hash([]byte(pw))
		if err != nil {
			t.Fatalf("HashPassword() error = %v", err)
		}

		if err = password.Check([]byte(pw), hash); err != nil {
			t.Fatalf("CheckPasswordHash() = false, want true")
		}
	})
}
