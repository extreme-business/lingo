// Package app represents the a set of functionality that the account domain provides.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/extreme-business/lingo/apps/account/auth/authentication"
	"github.com/extreme-business/lingo/apps/account/auth/registration"
	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/domain/user"
	"github.com/google/uuid"
)

// App is the application service for the account domain.
type App struct {
	logger              *slog.Logger
	userReader          *user.Reader
	authenticator       *authentication.Authenticator
	registrationManager *registration.Manager
}

type Config struct {
	Logger              *slog.Logger
	UserReader          *user.Reader
	Authenticator       *authentication.Authenticator
	RegistrationManager *registration.Manager
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if c.Logger == nil {
		return errors.New("logger is nil")
	}
	if c.UserReader == nil {
		return errors.New("user reader is nil")
	}
	if c.Authenticator == nil {
		return errors.New("authenticator is nil")
	}
	if c.RegistrationManager == nil {
		return errors.New("registration manager is nil")
	}
	return nil
}

func New(c Config) (*App, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return &App{
		logger:              c.Logger,
		userReader:          c.UserReader,
		authenticator:       c.Authenticator,
		registrationManager: c.RegistrationManager,
	}, nil
}

type RegisterUser struct {
	OrganizationID uuid.UUID
	DisplayName    string
	Email          string
	Password       string
}

func (r *App) RegisterUser(ctx context.Context, i RegisterUser) (*domain.User, error) {
	user, err := r.registrationManager.Register(ctx, registration.Registration{
		OrganizationID: i.OrganizationID,
		DisplayName:    i.DisplayName,
		Email:          i.Email,
		Password:       i.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	return user, nil
}

// ListUsers lists users.
// page starts from 0.
func (r *App) ListUsers(ctx context.Context, page uint) ([]*domain.User, error) {
	users, err := r.userReader.List(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
