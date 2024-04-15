package app

import (
	"context"
	"log/slog"

	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/dwethmar/lingo/cmd/auth/user/authentication"
	"github.com/dwethmar/lingo/cmd/auth/user/registration"
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

	user, err := r.registrationManager.Register(ctx, registration.Registration{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
