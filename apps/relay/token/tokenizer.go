package token

import (
	"time"

	"github.com/dwethmar/lingo/pkg/clock"
	jwt "github.com/golang-jwt/jwt/v5"
)

// Dispatcher sends a token to the user
type Tokenizer struct {
	clock     *clock.Clock
	secretKey []byte
	expiry    time.Duration
}

// Dispatch sends a token to the user
func (r *Tokenizer) Create(sub string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": sub,
			"exp": r.clock.Now().Add(r.expiry).Unix(),
		},
	)

	tokenString, err := token.SignedString(r.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
