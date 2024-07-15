package registration

import (
	"context"
	"errors"

	"github.com/extreme-business/lingo/cmd/account/domain"
	"github.com/extreme-business/lingo/cmd/account/password"
	"github.com/extreme-business/lingo/cmd/account/storage"
	"github.com/extreme-business/lingo/pkg/clock"
	"github.com/extreme-business/lingo/pkg/uuidgen"
)

// Manager is a manager for registration.
type Manager struct {
	uuidgen               uuidgen.Generator
	clock                 clock.Now
	userRepo              storage.UserRepository
	registrationValidator *registrationValidator
}

// Config is the configuration for the manager.
type Config struct {
	UUIDgen  uuidgen.Generator
	Clock    clock.Now
	UserRepo storage.UserRepository
}

// NewManager creates a new manager.
func NewManager(c Config) *Manager {
	return &Manager{
		uuidgen:               c.UUIDgen,
		clock:                 c.Clock,
		userRepo:              c.UserRepo,
		registrationValidator: newRegistrationValidator(),
	}
}

// configured checks if the manager is configured.
func (m *Manager) configured() error {
	if m.uuidgen == nil {
		return errors.New("uuidgen is nil")
	}
	if m.clock == nil {
		return errors.New("clock is nil")
	}
	if m.userRepo == nil {
		return errors.New("user repository is nil")
	}
	if m.registrationValidator == nil {
		return errors.New("registration validator is nil")
	}

	return nil

}

type Registration struct {
	User     *domain.User
	Password string
}

// CreateUser creates a new user.
func (m *Manager) Register(ctx context.Context, r Registration) (*domain.User, error) {
	if err := m.configured(); err != nil {
		return nil, err
	}

	if err := m.registrationValidator.Validate(r); err != nil {
		return nil, err
	}

	hashedPassword, err := password.Hash(r.Password)
	if err != nil {
		return nil, err
	}

	now := m.clock()
	s := &storage.User{}
	s.FromDomain(r.User)
	s.ID = m.uuidgen()
	s.HashedPassword = hashedPassword
	s.CreateTime = now
	s.UpdateTime = now

	user, err := m.userRepo.Create(ctx, s)
	if err != nil {
		return nil, err
	}

	user.HashedPassword = "" // Do not return the password
	var u domain.User
	user.ToDomain(&u)

	return &u, nil
}
