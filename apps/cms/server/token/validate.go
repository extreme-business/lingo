package token

import (
	"context"

	"github.com/extreme-business/lingo/pkg/httpmiddleware"
	"github.com/extreme-business/lingo/pkg/token"
)

var _ httpmiddleware.TokenValidator = &TokenValidator{}

type TokenValidator struct {
	tokenValidator *token.Validator
}

func NewTokenValidator(secret []byte) *TokenValidator {
	return &TokenValidator{
		tokenValidator: token.NewValidator(secret),
	}
}

// Validate implements httpmiddleware.Validator.
func (v *TokenValidator) Validate(_ context.Context, value string) error {
	_, err := v.tokenValidator.Validate(value)
	return err
}
