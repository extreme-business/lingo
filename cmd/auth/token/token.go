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

func ExtractExpirationTime(tokenString string) (time.Time, error) {
	var expirationTime time.Time
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if c, err := claims.GetExpirationTime(); err == nil {
			if c != nil {
				expirationTime = c.Time
			}
		} else {
			return time.Time{}, fmt.Errorf("failed to get expiration time: %w", err)
		}
	}

	if expirationTime.IsZero() {
		return time.Time{}, fmt.Errorf("invalid token payload")
	}

	return expirationTime, nil
}
