package password

import "testing"

func TestHashPassword(t *testing.T) {
	t.Run("should hash a password", func(t *testing.T) {
		password := "password"
		hash, err := Hash(password)
		if err != nil {
			t.Fatalf("HashPassword() error = %v", err)
		}

		if !Check(password, hash) {
			t.Fatalf("CheckPasswordHash() = false, want true")
		}
	})
}
