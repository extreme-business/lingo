package registration

import (
	"context"

	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/dwethmar/lingo/cmd/auth/password"
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/dwethmar/lingo/pkg/uuidgen"
)

// Manager is a manager for registration.
type Manager struct {
	uuidgen               *uuidgen.Generator
	clock                 *clock.Clock
	userRepo              user.Repository
	registrationValidator *registrationValidator
}

// Config is the configuration for the manager.
type Config struct {
	UUIDgen  *uuidgen.Generator
	Clock    *clock.Clock
	UserRepo user.Repository
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

type Registration struct {
	Username string
	Email    string
	Password string
}

// CreateUser creates a new user.
func (m *Manager) Register(ctx context.Context, r Registration) (*domain.User, error) {
	if err := m.registrationValidator.Validate(r); err != nil {
		return nil, err
	}

	hashedPassword, err := password.Hash(r.Password)
	if err != nil {
		return nil, err
	}

	user, err := m.userRepo.Create(ctx, &user.User{
		ID:         m.uuidgen.New(),
		Username:   r.Username,
		Email:      r.Email,
		Password:   hashedPassword,
		CreateTime: m.clock.Now(),
		UpdateTime: m.clock.Now(),
	})
	if err != nil {
		return nil, err
	}

	user.Password = "" // Do not return the password
	var u domain.User
	user.ToDomain(&u)

	return &u, nil
}
