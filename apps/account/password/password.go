package password

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultCost = 14
)

func Hash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), defaultCost)
}

func Check(password, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
