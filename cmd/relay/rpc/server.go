package rpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dwethmar/lingo/cmd/relay/token"
	relayProto "github.com/dwethmar/lingo/proto/v1/relay"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	relayProto.UnimplementedRelayServiceServer
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

func (s *Server) CreateRegistrationToken(ctx context.Context, req *relayProto.RegistrationTokenRequest) (*emptypb.Empty, error) {
	s.logger.Info("CreateRegistrationToken")

	if err := s.RegistrationTokenManager.Create(req.Email); err != nil {
		return nil, fmt.Errorf("failed to send token: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) CreateMessage(ctx context.Context, req *relayProto.CreateMessageRequest) (*emptypb.Empty, error) {
	s.logger.Info("CreateMessage")

	return &emptypb.Empty{}, nil
}
