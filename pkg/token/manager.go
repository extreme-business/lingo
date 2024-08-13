package token

import (
	"fmt"
	"time"
)

// Created is the event that is dispatched when a token is created.
type Created struct {
	Email string
	Token string // hex encoded
}

// Manager manages the token validation and dispatching.
type Manager struct {
	Validator *Validator
	Tokenizer *Tokenizer
}

func NewManager(
	clock func() time.Time,
	signingKey []byte,
	tokenExpireDuration time.Duration,
) *Manager {
	return &Manager{
		Validator: &Validator{secretKey: signingKey},
		Tokenizer: &Tokenizer{clock: clock, secretKey: signingKey, expiry: tokenExpireDuration},
	}
}

func (m *Manager) Create(id string) (string, error) {
	token, err := m.Tokenizer.Create(id)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Validate validates the token and returns the email hash.
// Token is hex encoded.
func (m *Manager) Validate(token string) (*Claims, error) {
	claims, err := m.Validator.Validate(token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return claims, nil
}
