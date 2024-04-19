package authentication

import (
	"github.com/dwethmar/lingo/pkg/validate"
)

const (
	// minEmailLength is the minimum length of an email.
	minEmailLength = 1
	// maxEmailLength is the maximum length of an email.
	maxEmailLength = 50
	// minPasswordLength is the minimum length of a password.
	minPasswordLength = 1
	// maxPasswordLength is the maximum length of a password.
	maxPasswordLength = 100
)

type credentialsValidator struct {
	emailValidator    validate.StringValidator
	passwordValidator validate.StringValidator
}

func newCredentialsValidator() *credentialsValidator {
	return &credentialsValidator{
		emailValidator: validate.StringValidator{
			validate.MinLength("email", minEmailLength),
			validate.MaxLength("email", maxEmailLength),
		},
		passwordValidator: validate.StringValidator{
			validate.MinLength("password", minPasswordLength),
			validate.MaxLength("password", maxPasswordLength),
		},
	}
}

func (v *credentialsValidator) Validate(r Credentials) error {
	if err := v.emailValidator.Validate(r.Email); err != nil {
		return err
	}

	if err := v.passwordValidator.Validate(r.Password); err != nil {
		return err
	}

	return nil
}
