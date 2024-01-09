package rpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dwethmar/lingo/cmd/relay/token"
	"github.com/dwethmar/lingo/gen/go/proto/relay/v1"
)

type Server struct {
	relay.UnimplementedRelayServiceServer
	logger                     *slog.Logger
	RegistrationTokenManager   *token.Manager
	AuthenticationTokenManager *token.Manager
}

func New(
	logger *slog.Logger,
	registrationTokenManager *token.Manager,
	authenticationTokenManager *token.Manager,
) *Server {
	return &Server{
		logger:                     logger,
		RegistrationTokenManager:   registrationTokenManager,
		AuthenticationTokenManager: authenticationTokenManager,
	}
}

func (s *Server) CreateRegisterToken(ctx context.Context, req *relay.CreateRegisterTokenRequest) (*relay.CreateRegisterTokenResponse, error) {
	s.logger.Info("CreateRegistrationToken")

	if err := s.RegistrationTokenManager.Create(req.Email); err != nil {
		return nil, fmt.Errorf("failed to send token: %w", err)
	}

	return &relay.CreateRegisterTokenResponse{}, nil
}

func (s *Server) CreateMessage(ctx context.Context, req *relay.CreateMessageRequest) (*relay.CreateMessageResponse, error) {
	s.logger.Info("CreateMessage")

	return &relay.CreateMessageResponse{}, nil
}
