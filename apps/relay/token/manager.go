package token

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dwethmar/lingo/pkg/clock"
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
	clock *clock.Clock,
	signingKey []byte,
	tokenExpireDuration time.Duration,
	dispatchCh chan<- Created,
) *Manager {
	return &Manager{
		Validator:  &Validator{secretKey: signingKey},
		Tokenizer:  &Tokenizer{clock: clock, secretKey: signingKey, expiry: tokenExpireDuration},
		dispatchCh: dispatchCh,
	}
}

func (m *Manager) Create(email string) error {
	h := sha256.New()
	h.Write([]byte(email))
	emailhash := h.Sum(nil)

	token, err := m.Tokenizer.Create(string(emailhash))
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
func (m *Manager) Validate(token string) (Claims, error) {
	claims, err := m.Validator.Validate(token)
	if err != nil {
		return Claims{}, fmt.Errorf("failed to validate token: %w", err)
	}

	return claims, nil
}
