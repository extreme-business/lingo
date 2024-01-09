package token

import (
	"crypto/sha256"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Dispatcher sends a token to the user
type Tokenizer struct {
	secretKey []byte
	expiry    time.Duration
}

// Dispatch sends a token to the user
func (r *Tokenizer) Create(email string) (string, error) {
	h := sha256.New()
	h.Write([]byte(email))
	emailhash := h.Sum(nil)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": string(emailhash),
			"exp": time.Now().Add(r.expiry).Unix(),
		},
	)

	tokenString, err := token.SignedString(r.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
