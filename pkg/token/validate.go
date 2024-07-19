package token

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Validator validates a token.
type Validator struct {
	secretKey []byte
}

func NewValidator(secretKey []byte) *Validator {
	return &Validator{
		secretKey: secretKey,
	}
}

type Claims struct {
	ExpirationTime time.Time
	Sub            string
}

// Validate validates the token and returns the email hash.
func (v *Validator) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return v.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}

		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("sub is not a string: %w", ErrInvalidTokenClaims)
	}

	var expirationTime time.Time
	numericDate, err := claims.GetExpirationTime()
	if err == nil {
		if numericDate != nil {
			expirationTime = numericDate.Time
		}
	} else {
		return nil, fmt.Errorf("expiration time is not valid: %w", ErrInvalidTokenClaims)
	}

	return &Claims{
		ExpirationTime: expirationTime,
		Sub:            sub,
	}, nil
}
