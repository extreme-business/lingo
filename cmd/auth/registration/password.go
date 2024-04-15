package registration

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password is too short (min 8 characters)")
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword checks if the password is valid
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}
