package registration

import (
	"context"
	"errors"
	"time"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/password"
	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/extreme-business/lingo/pkg/uuidgen"
)

// Manager is a manager for registration.
type Manager struct {
	uuidgen               uuidgen.Generator
	clock                 func() time.Time
	userRepo              storage.UserRepository
	registrationValidator *registrationValidator
}

// Config is the configuration for the manager.
type Config struct {
	UUIDgen  uuidgen.Generator
	Clock    func() time.Time
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

	hashedPassword, err := password.Hash([]byte(r.Password))
	if err != nil {
		return nil, err
	}

	now := m.clock()
	s := &storage.User{}
	if err := r.User.ToStorage(s); err != nil {
		return nil, err
	}

	s.ID = m.uuidgen()
	s.HashedPassword = string(hashedPassword)
	s.CreateTime = now
	s.UpdateTime = now

	user, err := m.userRepo.Create(ctx, s)
	if err != nil {
		return nil, err
	}

	user.HashedPassword = "" // Do not return the password

	u := &domain.User{}
	return u, u.FromStorage(user)
}
