package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	defaultCost = 14
)

var (
	ErrMismatchedHashAndPassword = errors.New("hashedPassword is not the hash of the given password")
)

func Hash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, defaultCost)
}

func Check(password, hash []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrMismatchedHashAndPassword
	}

	return err
}
