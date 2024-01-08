package register

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

var (
	TokenExpirationDuration = time.Minute * 15
)

type Token struct {
	EmailHash       []byte `json:"email_hash"`
	CreateTimestamp int64  `json:"create_timestamp"`
}

type Handler interface {
	SendToken(email string, token []byte) error
}

// Registrar is a service that creates accounts
type Registrar struct {
	cipher Cipher
	h      Handler
}

// New creates a new Registrar
func New(cipher Cipher, h Handler) *Registrar {
	return &Registrar{
		cipher: cipher,
		h:      h,
	}
}

// ValidateToken validates a token and returns the email hash
func (r *Registrar) ValidateToken(token []byte) ([]byte, error) {
	decrypted, err := r.cipher.Decrypt(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	t := &Token{}
	if err := json.Unmarshal(decrypted, t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	if t.CreateTimestamp < time.Now().Add(-TokenExpirationDuration).UnixMilli() {
		return nil, fmt.Errorf("token expired")
	}

	return t.EmailHash, nil
}

// Register creates a new account
func (r *Registrar) SendToken(email string) error {
	h := sha256.New()
	h.Write([]byte(email))

	t := &Token{
		EmailHash:       h.Sum(nil),
		CreateTimestamp: time.Now().UnixMilli(),
	}

	json, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	encrypted, err := r.cipher.Encrypt(json)
	if err != nil {
		return fmt.Errorf("failed to encrypt: %w", err)
	}

	if err := r.h.SendToken(email, encrypted); err != nil {
		return fmt.Errorf("failed to send token: %w", err)
	}

	return nil
}
