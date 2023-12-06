package authentication

import (
	"context"
	"errors"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/dwethmar/lingo/cmd/auth/password"
	"github.com/dwethmar/lingo/cmd/auth/storage"
	"github.com/dwethmar/lingo/cmd/auth/token"
	"github.com/dwethmar/lingo/pkg/clock"
)

const (
	authTokenDuration    = 5 * time.Minute
	refreshTokenDuration = 30 * time.Minute
)

type Manager struct {
	credentialsValidator *credentialsValidator
	userRepo             storage.UserRepository
	AuthTokenManager     *token.Manager
	RefreshTokenManager  *token.Manager
}

type Config struct {
	Clock                    clock.Now
	SigningKeyRegistration   []byte
	SigningKeyAuthentication []byte
	UserRepo                 storage.UserRepository
}

func NewManager(c Config) *Manager {
	return &Manager{
		credentialsValidator: newCredentialsValidator(),
		userRepo:             c.UserRepo,
		AuthTokenManager: token.NewManager(
			c.Clock,
			c.SigningKeyRegistration,
			authTokenDuration,
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
	Token        string
	RefreshToken string
}

// Authenticate authenticates a user with the given credentials.
func (m *Manager) Authenticate(ctx context.Context, c Credentials) (*Authentication, error) {
	if err := m.credentialsValidator.Validate(c); err != nil {
		return nil, err
	}

	u, err := m.userRepo.GetByEmail(ctx, c.Email)
	if err != nil {
		return nil, err
	}

	if !password.Check(c.Password, u.Password) {
		return nil, errors.New("could not authenticate")
	}

	authToken, err := m.AuthTokenManager.New(u.ID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.RefreshTokenManager.New(u.ID.String())
	if err != nil {
		return nil, err
	}

	var domainUser domain.User
	u.ToDomain(&domainUser)

	return &Authentication{
		User:         &domainUser,
		Token:        authToken,
		RefreshToken: refreshToken,
	}, nil
}
