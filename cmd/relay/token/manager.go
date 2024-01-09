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
	Validator  *Validator
	Tokenizer  *Tokenizer
	dispatchCh chan<- Created
}

func NewManager(
	signingKey []byte,
	tokenExpireDuration time.Duration,
	dispatchCh chan<- Created,
) *Manager {
	return &Manager{
		Validator:  &Validator{secretKey: signingKey},
		Tokenizer:  &Tokenizer{secretKey: signingKey, expiry: tokenExpireDuration},
		dispatchCh: dispatchCh,
	}
}

func (m *Manager) Create(email string) error {
	token, err := m.Tokenizer.Create(email)
	if err != nil {
		return err
	}

	m.dispatchCh <- Created{
		Email: email,
		Token: token,
	}

	return nil
}

// Validate validates the token and returns the email hash.
// Token is hex encoded.
func (m *Manager) Validate(token string) (string, error) {
	emailHash, err := m.Validator.Validate(token)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	return emailHash, nil
}
