package registration

import (
	"github.com/dwethmar/lingo/pkg/validate"
)

const (
	// maxUsernameLength is the maximum length of a username.
	maxUsernameLength = 50
	// minUsernameLength is the minimum length of a username.
	minUsernameLength = 3
	// maxEmailLength is the maximum length of an email.
	maxEmailLength = 50
	// minEmailLength is the minimum length of an email.
	minEmailLength = 3
	// maxPasswordLength is the maximum length of a password.
	maxPasswordLength = 50
	// minPasswordLength is the minimum length of a password.
	minPasswordLength = 8
	// minPasswordSpecialChars is the minimum amount of special characters in a password.
	minPasswordSpecialChars = 1
	// minPasswordDigits is the minimum amount of digits in a password.
	minPasswordDigits = 1
)

// usernameSpecialChars returns a list of special characters that are allowed in a username.
func usernameSpecialChars() []rune { return []rune{'_', '-'} }

type registrationValidator struct {
	usernameValidator validate.StringValidator
	emailValidator    validate.StringValidator
	passwordValidator validate.StringValidator
}

func newRegistrationValidator() *registrationValidator {
	return &registrationValidator{
		usernameValidator: validate.StringValidator{
			validate.MaxLength("username", maxUsernameLength),
			validate.MinLength("username", minUsernameLength),
			validate.SpecialCharWhitelist("username", usernameSpecialChars()...),
		},
		emailValidator: validate.StringValidator{
			validate.MaxLength("email", maxEmailLength),
			validate.MinLength("email", minEmailLength),
		},
		passwordValidator: validate.StringValidator{
			validate.MinLength("password", minPasswordLength),
			validate.MaxLength("password", maxPasswordLength),
			validate.ContainsSpecialChars("password", minPasswordSpecialChars),
			validate.ContainsDigits("password", minPasswordDigits),
		},
	}
}

func (v *registrationValidator) Validate(r Registration) error {
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
