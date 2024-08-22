package app

import (
	"context"
	"log/slog"

	"github.com/extreme-business/lingo/apps/account/auth/authentication"
	"github.com/extreme-business/lingo/apps/account/auth/registration"
	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/google/uuid"
)

// App is the application service for the account domain.
type App struct {
	logger              *slog.Logger
	authenticator       *authentication.Authenticator
	registrationManager *registration.Manager
}

func New(
	logger *slog.Logger,
	authenticator *authentication.Authenticator,
	registrationManager *registration.Manager,
) *App {
	return &App{
		logger:              logger,
		authenticator:       authenticator,
		registrationManager: registrationManager,
	}
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
		return nil, err
	}
	return user, nil
}
