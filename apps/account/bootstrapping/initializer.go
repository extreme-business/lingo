// Package bootstrapping contains the Initializer to set up the system data.
package bootstrapping

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/extreme-business/lingo/cmd/account/password"
	"github.com/extreme-business/lingo/cmd/account/storage"
	"github.com/extreme-business/lingo/pkg/clock"
	"github.com/extreme-business/lingo/pkg/validate"
	"github.com/google/uuid"
)

const (
	systemUserName = "system"
)

type SystemUserConfig struct {
	ID       uuid.UUID // ID is the user id of the system user.
	Email    string    // Email is the email of the system user.
	Password string    // Password is the password of the system user.
}

func (c *SystemUserConfig) Validate() error {
	idValidator := validate.UUIDValidator{validate.UUIDIsNotNil("ID")}
	emailValidator := validate.StringValidator{validate.StringMinLength("Email", 1)}
	passwordValidator := validate.StringValidator{validate.StringMinLength("Password", 1)}

	if err := idValidator.Validate(c.ID); err != nil {
		return err
	}

	if err := emailValidator.Validate(c.Email); err != nil {
		return err
	}

	if err := passwordValidator.Validate(c.Password); err != nil {
		return err
	}

	return nil
}

// SystemOrgConfig is the configuration for the system organization.
type SystemOrgConfig struct {
	ID        uuid.UUID // ID is the organization id of the system user.
	LegalName string    // LegalName is the organization name of the system user.
	Slug      string    // Slug is the organization slug of the system user.
}

// systemOrganizationConfigValidator returns a function that validates the system organization configuration.
func (c *SystemOrgConfig) Validate() error {
	idValidator := validate.UUIDValidator{
		validate.UUIDIsNotNil("ID"),
	}
	legalNameValidator := validate.StringValidator{
		validate.StringIsUtf8("LegalName"),
		validate.StringRequired("LegalName"),
		validate.StringMinLength("LegalName", 1),
	}
	slugValidator := validate.StringValidator{
		validate.StringIsUtf8("Slug"),
		validate.StringRequired("Slug"),
		validate.StringMinLength("Slug", 1),
		validate.StringMaxLength("Slug", 100),
		validate.SpecialCharWhitelist("Slug", '-'),
	}

	if err := idValidator.Validate(c.ID); err != nil {
		return err
	}

	if err := legalNameValidator.Validate(c.LegalName); err != nil {
		return err
	}

	if err := slugValidator.Validate(c.Slug); err != nil {
		return err
	}

	return nil
}

// Initializer is responsible for setting up the system user and organization.
type Initializer struct {
	logger                   *slog.Logger
	systemUserConfig         SystemUserConfig
	systemOrganizationConfig SystemOrgConfig
	clock                    clock.Now
	dbManager                storage.DBManager
}

type Config struct {
	Logger                   *slog.Logger
	SystemUserConfig         SystemUserConfig
	SystemOrganizationConfig SystemOrgConfig
	Clock                    clock.Now
	DBManager                storage.DBManager
}

func (c Config) Validate() error {
	if c.Logger == nil {
		return errors.New("logger is required")
	}

	if c.Clock == nil {
		return errors.New("clock is required")
	}

	if c.DBManager == nil {
		return errors.New("db manager is required")
	}

	if err := c.SystemUserConfig.Validate(); err != nil {
		return fmt.Errorf("invalid system user config: %w", err)
	}

	if err := c.SystemOrganizationConfig.Validate(); err != nil {
		return fmt.Errorf("invalid system organization config: %w", err)
	}

	return nil
}

func New(config Config) (*Initializer, error) {
	return &Initializer{
		logger:                   config.Logger,
		systemUserConfig:         config.SystemUserConfig,
		systemOrganizationConfig: config.SystemOrganizationConfig,
		clock:                    config.Clock,
		dbManager:                config.DBManager,
	}, config.Validate()
}

// Setup sets up the system user and organization.
func (s *Initializer) Setup(ctx context.Context) error {
	return s.dbManager.BeginOp(ctx, func(ctx context.Context, r storage.Repositories) error {
		if r.User == nil {
			return errors.New("user repository is required")
		}

		if r.Organization == nil {
			return errors.New("organization repository is required")
		}

		return s.setup(ctx, &r)
	})
}

// setupOrganization sets up the system organization. If the organization already exists, it will be updated if necessary.
func (s *Initializer) setupOrganization(ctx context.Context, r storage.OrganizationRepository) (*storage.Organization, error) {
	now := s.clock()

	// check if the organization already exists
	org, err := r.Get(ctx, s.systemOrganizationConfig.ID)
	if err == nil {
		// check if the organization needs to be updated
		updated := false
		triggeredChecks := []string{}
		for check, isDifferent := range map[string]func(*storage.Organization) bool{
			"legal name": func(o *storage.Organization) bool { return o.LegalName != s.systemOrganizationConfig.LegalName },
		} {
			if isDifferent(org) {
				triggeredChecks = append(triggeredChecks, check)
				updated = true
				break
			}
		}

		if updated {
			s.logger.Info("system organization update triggered", slog.String("changes", strings.Join(triggeredChecks, ", ")))

			org.LegalName = s.systemOrganizationConfig.LegalName
			org.UpdateTime = now

			o, uErr := r.Update(ctx, org, []storage.OrganizationField{
				storage.OrganizationLegalName,
			})

			if uErr != nil {
				return nil, fmt.Errorf("failed to update system organization: %w", uErr)
			}

			return o, nil
		}
	}

	// if the organization does not exist, create it
	if errors.Is(err, storage.ErrOrganizationNotFound) {
		s.logger.Info("system organization creation triggered")

		o, cErr := r.Create(ctx, &storage.Organization{
			ID:         s.systemOrganizationConfig.ID,
			LegalName:  s.systemOrganizationConfig.LegalName,
			UpdateTime: now,
			CreateTime: now,
		})

		if cErr != nil {
			return nil, fmt.Errorf("failed to create system organization: %w", cErr)
		}

		return o, nil
	}

	return org, err
}

// setupUser sets up the system user. If the user already exists, it will be updated if necessary.
func (s *Initializer) setupUser(ctx context.Context, org *storage.Organization, r storage.UserRepository) (*storage.User, error) {
	now := s.clock()

	hashedPassword, hErr := password.Hash(s.systemUserConfig.Password)
	if hErr != nil {
		return nil, hErr
	}

	// check if the user already exists
	user, err := r.Get(ctx, s.systemUserConfig.ID)
	if err == nil {
		// check if the user needs to be updated
		updated := false
		triggeredChecks := []string{}

		for check, isDifferent := range map[string]func(*storage.User) bool{
			"organization id": func(u *storage.User) bool { return u.OrganizationID != org.ID },
			"display name":    func(u *storage.User) bool { return u.DisplayName != systemUserName },
			"email":           func(u *storage.User) bool { return u.Email != s.systemUserConfig.Email },
			"password":        func(u *storage.User) bool { return !password.Check(u.HashedPassword, hashedPassword) },
		} {
			if isDifferent(user) {
				updated = true
				triggeredChecks = append(triggeredChecks, check)
			}
		}

		if updated {
			s.logger.Info("system user update triggered", slog.String("changes", strings.Join(triggeredChecks, ", ")))

			user.OrganizationID = org.ID
			user.DisplayName = systemUserName
			user.Email = s.systemUserConfig.Email
			user.UpdateTime = now
			user.HashedPassword = hashedPassword

			u, uErr := r.Update(ctx, user, []storage.UserField{
				storage.UserOrganizationID,
				storage.UserDisplayName,
				storage.UserEmail,
				storage.UserPassword,
			})

			if uErr != nil {
				return nil, fmt.Errorf("failed to update system user: %w", uErr)
			}

			return u, nil
		}
	}

	// if the user does not exist, create it
	if errors.Is(err, storage.ErrUserNotFound) {
		s.logger.Info("system user creation triggered")

		u, cErr := r.Create(ctx, &storage.User{
			ID:             s.systemUserConfig.ID,
			OrganizationID: org.ID,
			DisplayName:    systemUserName,
			Email:          s.systemUserConfig.Email,
			HashedPassword: hashedPassword,
			UpdateTime:     now,
			CreateTime:     now,
		})

		if cErr != nil {
			return nil, fmt.Errorf("failed to create system user: %w", cErr)
		}

		return u, nil
	}

	return nil, err
}

func (s *Initializer) setup(ctx context.Context, r *storage.Repositories) error {
	if s.clock == nil {
		return errors.New("clock is required")
	}

	// Create the system organization and user.
	org, err := s.setupOrganization(ctx, r.Organization)
	if err != nil {
		return fmt.Errorf("failed to set up system organization: %w", err)
	}

	_, err = s.setupUser(ctx, org, r.User)
	if err != nil {
		return fmt.Errorf("failed to set up system user: %w", err)
	}

	return nil
}
