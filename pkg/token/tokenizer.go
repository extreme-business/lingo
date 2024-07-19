package token

import (
	"time"

	"github.com/extreme-business/lingo/pkg/clock"
	"github.com/golang-jwt/jwt/v5"
)

// Tokenizer creates a token.
type Tokenizer struct {
	clock     clock.Now
	secretKey []byte
	expiry    time.Duration
}

// Dispatch sends a token to the user.
func (r *Tokenizer) Create(sub string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": sub,
			"exp": r.clock().Add(r.expiry).Unix(),
		},
	)

	tokenString, err := token.SignedString(r.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
