package registration

import (
	"github.com/dwethmar/lingo/pkg/validate"
)

type RegistrationValidator struct {
	usernameValidator validate.StringValidator
	emailValidator    validate.StringValidator
	passwordValidator validate.StringValidator
}

func NewRegistrationValidator() *RegistrationValidator {
	return &RegistrationValidator{
		usernameValidator: validate.StringValidator{
			validate.MaxLength("username", 50),
			validate.MinLength("username", 3),
			validate.SpecialCharWhitelist("username", '_', '-'),
		},
		emailValidator: validate.StringValidator{
			validate.MaxLength("email", 50),
			validate.MinLength("email", 3),
		},
		passwordValidator: validate.StringValidator{
			validate.MinLength("password", 8),
			validate.MaxLength("password", 50),
			validate.ContainsSpecialChars("password", 1),
			validate.ContainsDigits("password", 1),
		},
	}
}

func (v *RegistrationValidator) Validate(r Registration) error {
	if err := v.usernameValidator.Validate(r.Username); err != nil {
		return err
	}

	if err := v.emailValidator.Validate(r.Email); err != nil {
		return err
	}

	if err := v.passwordValidator.Validate(r.Password); err != nil {
		return err
	}

	return nil
}
