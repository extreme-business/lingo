package server

import (
	"context"

	"github.com/dwethmar/lingo/apps/relay"
	protorelay "github.com/dwethmar/lingo/protogen/go/proto/private/relay/v1"
)

type Service struct {
	protorelay.UnimplementedRelayServiceServer
	relay *relay.Relay
}

func New(relay *relay.Relay) *Service {
	return &Service{
		relay: relay,
	}
}

func (s *Service) CreateRegisterToken(ctx context.Context, req *protorelay.CreateRegisterTokenRequest) (*protorelay.CreateRegisterTokenResponse, error) {
	err := s.relay.CreateRegisterToken(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return &protorelay.CreateRegisterTokenResponse{}, nil
}

func (s *Service) CreateMessage(ctx context.Context, req *protorelay.CreateMessageRequest) (*protorelay.CreateMessageResponse, error) {
	err := s.relay.CreateMessage(ctx, req.Message)
	if err != nil {
		return nil, err
	}

	return &protorelay.CreateMessageResponse{}, nil
}
