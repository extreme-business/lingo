package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	defaultCost = 14
)

func Hash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, defaultCost)
}

func Check(password, hash []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return fmt.Errorf("password does not match hash: %w", err)
	}
	return nil
}
