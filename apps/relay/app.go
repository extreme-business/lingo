package relay

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dwethmar/lingo/apps/relay/token"
)

type Relay struct {
	logger                     *slog.Logger
	RegistrationTokenManager   *token.Manager
	AuthenticationTokenManager *token.Manager
}

func New(
	logger *slog.Logger,
	registrationTokenManager *token.Manager,
	authenticationTokenManager *token.Manager,
) *Relay {
	return &Relay{
		logger:                     logger,
		RegistrationTokenManager:   registrationTokenManager,
		AuthenticationTokenManager: authenticationTokenManager,
	}
}

func (r *Relay) CreateRegisterToken(ctx context.Context, email string) error {
	r.logger.Info("CreateRegistrationToken")

	if err := r.RegistrationTokenManager.Create(email); err != nil {
		return fmt.Errorf("failed to send token: %w", err)
	}

	return nil
}

func (r *Relay) CreateMessage(ctx context.Context, message string) error {
	r.logger.Info("CreateMessage")

	return nil
}
