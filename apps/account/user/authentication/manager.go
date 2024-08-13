package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/password"
	"github.com/extreme-business/lingo/apps/account/storage"
	"github.com/extreme-business/lingo/pkg/token"
)

const (
	accountTokenDuration = 5 * time.Minute
	refreshTokenDuration = 30 * time.Minute
)

type Manager struct {
	credentialsValidator *credentialsValidator
	userRepo             storage.UserRepository
	AccountTokenManager  *token.Manager
	RefreshTokenManager  *token.Manager
}

type Config struct {
	Clock                    func() time.Time
	SigningKeyRegistration   []byte
	SigningKeyAuthentication []byte
	UserRepo                 storage.UserRepository
}

func NewManager(c Config) *Manager {
	return &Manager{
		credentialsValidator: newCredentialsValidator(),
		userRepo:             c.UserRepo,
		AccountTokenManager: token.NewManager(
			c.Clock,
			c.SigningKeyRegistration,
			accountTokenDuration,
		),
		RefreshTokenManager: token.NewManager(
			c.Clock,
			c.SigningKeyAuthentication,
			refreshTokenDuration,
		),
	}
}

type Credentials struct {
	Email    string
	Password string
}

// Authentication is the process of verifying whether someone is who they claim to be when accessing a system.
type Authentication struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
}

// Authenticate a user with the given credentials.
func (m *Manager) Authenticate(ctx context.Context, c Credentials) (*Authentication, error) {
	if err := m.credentialsValidator.Validate(c); err != nil {
		return nil, err
	}

	u, err := m.userRepo.GetByEmail(ctx, c.Email)
	if err != nil {
		return nil, err
	}

	if err = password.Check([]byte(c.Password), []byte(u.HashedPassword)); err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	accountToken, err := m.AccountTokenManager.Create(u.ID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.RefreshTokenManager.Create(u.ID.String())
	if err != nil {
		return nil, err
	}

	var domainUser domain.User
	if err := domainUser.FromStorage(u); err != nil {
		return nil, err
	}

	return &Authentication{
		User:         &domainUser,
		AccessToken:  accountToken,
		RefreshToken: refreshToken,
	}, nil
}
