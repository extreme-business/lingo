package token

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenMalformed     = errors.New("token is malformed")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid claims")
)

// Validator validates a token.
type Validator struct {
	secretKey []byte
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
	if numericDate, err := claims.GetExpirationTime(); err == nil {
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
