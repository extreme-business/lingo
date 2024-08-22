package registration

import (
	"context"
	"errors"
	"fmt"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/domain/user"
	"github.com/extreme-business/lingo/apps/account/password"
	"github.com/extreme-business/lingo/pkg/uuidgen"
	"github.com/google/uuid"
)

// Manager is a manager for registration.
type Manager struct {
	genUUID               uuidgen.Generator
	userWriter            *user.Writer
	registrationValidator *registrationValidator
}

// Config is the configuration for the manager.
type Config struct {
	GenUUID    uuidgen.Generator
	UserWriter *user.Writer
}

// NewManager creates a new manager.
func NewManager(c Config) *Manager {
	return &Manager{
		genUUID:               c.GenUUID,
		userWriter:            c.UserWriter,
		registrationValidator: newRegistrationValidator(),
	}
}

// configured checks if the manager is configured.
func (m *Manager) configured() error {
	if m.genUUID == nil {
		return errors.New("uuidgen is nil")
	}
	if m.userWriter == nil {
		return errors.New("user repository is nil")
	}
	if m.registrationValidator == nil {
		return errors.New("registration validator is nil")
	}

	return nil
}

type Registration struct {
	OrganizationID uuid.UUID
	DisplayName    string
	Email          string
	Password       string
}

// CreateUser creates a new user.
func (m *Manager) Register(ctx context.Context, r Registration) (*domain.User, error) {
	if err := m.configured(); err != nil {
		return nil, fmt.Errorf("manager not configured: %w", err)
	}
	if err := m.registrationValidator.Validate(r); err != nil {
		return nil, err
	}
	hashedPassword, err := password.Hash([]byte(r.Password))
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}
	u := &domain.User{
		ID:             m.genUUID(),
		Email:          r.Email,
		DisplayName:    r.DisplayName,
		OrganizationID: r.OrganizationID,
		HashedPassword: string(hashedPassword),
	}
	user, err := m.userWriter.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	user.HashedPassword = "" // Do not return the password
	return user, nil
}
