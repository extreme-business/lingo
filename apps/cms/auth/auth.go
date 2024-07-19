package auth

import (
	"context"

	accountproto "github.com/extreme-business/lingo/proto/gen/go/public/account/v1"
)

type Authenticator struct {
	client accountproto.AccountServiceClient
}

type SuccessResponse struct {
	AccessToken  string
	RefreshToken string
}

func (a *Authenticator) Authenticate(ctx context.Context, email, password string) (*SuccessResponse, error) {
	r, err := a.client.LoginUser(ctx, &accountproto.LoginUserRequest{
		Login: &accountproto.LoginUserRequest_Email{
			Email: email,
		},
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResponse{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
	}, nil
}

func NewAuthenticator(client accountproto.AccountServiceClient) *Authenticator {
	return &Authenticator{
		client: client,
	}
}
