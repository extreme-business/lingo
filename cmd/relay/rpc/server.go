package rpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dwethmar/lingo/cmd/relay/register"
	relayProto "github.com/dwethmar/lingo/proto/v1/relay"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	relayProto.UnimplementedRelayServiceServer
	logger   *slog.Logger
	register *register.Registrar
}

func New(logger *slog.Logger, register *register.Registrar) *Server {
	return &Server{
		logger:   logger,
		register: register,
	}
}

func (s *Server) CreateRegistrationToken(ctx context.Context, req *relayProto.RegistrationTokenRequest) (*emptypb.Empty, error) {
	s.logger.Info("CreateRegistrationToken")

	if err := s.register.SendToken(req.Email); err != nil {
		return nil, fmt.Errorf("failed to send token: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) CreateMessage(ctx context.Context, req *relayProto.CreateMessageRequest) (*emptypb.Empty, error) {
	s.logger.Info("CreateMessage")

	return &emptypb.Empty{}, nil
}
