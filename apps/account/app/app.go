package app

import (
	"context"
	"log/slog"

	"github.com/extreme-business/lingo/apps/account/auth/authentication"
	"github.com/extreme-business/lingo/apps/account/auth/registration"
	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/google/uuid"
)

type Account struct {
	logger              *slog.Logger
	accountManager      *authentication.Manager
	registrationManager *registration.Manager
}

func New(
	logger *slog.Logger,
	accountManager *authentication.Manager,
	registrationManager *registration.Manager,
) *Account {
	return &Account{
		logger:              logger,
		accountManager:      accountManager,
		registrationManager: registrationManager,
	}
}

type RegisterUser struct {
	OrganizationID uuid.UUID
	DisplayName    string
	Email          string
	Password       string
}

func (r *Account) RegisterUser(ctx context.Context, i RegisterUser) (*domain.User, error) {
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
