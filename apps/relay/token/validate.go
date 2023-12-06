package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid claims")
)

type Validator struct {
	secretKey []byte
}

type Claims struct {
	ExpirationTime time.Time
	Sub            string
}

// Validate validates the token and returns the email hash.
func (v *Validator) Validate(tokenStr string) (Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return v.secretKey, nil
	})

	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, ErrInvalidTokenClaims
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return Claims{}, fmt.Errorf("sub is not a string: %w", ErrInvalidTokenClaims)
	}

	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return Claims{}, fmt.Errorf("failed to get expiration time: %w", err)
	}

	return Claims{
		ExpirationTime: expirationTime.Time,
		Sub:            sub,
	}, nil
}
