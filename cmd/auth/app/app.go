package app

import (
	"context"
	"log/slog"

	"github.com/dwethmar/lingo/cmd/auth/authentication"
	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/dwethmar/lingo/cmd/auth/registration"
)

type Auth struct {
	logger                *slog.Logger
	authenticationManager *authentication.Manager
	registrationManager   *registration.Manager
}

func New(
	logger *slog.Logger,
	authenticationManager *authentication.Manager,
	registrationManager *registration.Manager,
) *Auth {
	return &Auth{
		logger:                logger,
		authenticationManager: authenticationManager,
		registrationManager:   registrationManager,
	}
}

func (r *Auth) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	r.logger.Info("Register")

	return &domain.User{}, nil
}
