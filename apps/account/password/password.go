package password

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultCost = 14
)

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	return string(bytes), err
}

func Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
