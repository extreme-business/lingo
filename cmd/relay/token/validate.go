package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	secretKey []byte
}

// Validate validates the token and returns the email hash.
func (v *Validator) Validate(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return v.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	emailHash, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	return emailHash, nil
}
