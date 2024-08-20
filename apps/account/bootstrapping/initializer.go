// Package bootstrapping contains the Initializer to set up the system data.
package bootstrapping

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/extreme-business/lingo/apps/account/domain"
	"github.com/extreme-business/lingo/apps/account/domain/organization"
	"github.com/extreme-business/lingo/apps/account/domain/user"
	"github.com/extreme-business/lingo/apps/account/password"
	"github.com/extreme-business/lingo/apps/account/storage"
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
		validate.StringUtf8("LegalName"),
		validate.StringNotEmpty("LegalName"),
		validate.StringMinLength("LegalName", 1),
	}
	slugValidator := validate.StringValidator{
		validate.StringUtf8("Slug"),
		validate.StringNotEmpty("Slug"),
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
	clock                    func() time.Time
	dbManager                storage.DBManager
}

type Config struct {
	Logger                   *slog.Logger
	SystemUserConfig         SystemUserConfig
	SystemOrganizationConfig SystemOrgConfig
	Clock                    func() time.Time
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
func (s *Initializer) setupOrganization(ctx context.Context, r *organization.Reader, w *organization.Writer) (*domain.Organization, error) {
	now := s.clock()

	// check if the organization already exists
	org, err := r.Get(ctx, s.systemOrganizationConfig.ID)
	if err == nil {
		// check if the organization needs to be updated
		changes := []storage.OrganizationField{}
		for check, diff := range map[storage.OrganizationField]func(*domain.Organization) bool{
			storage.OrganizationLegalName: func(o *domain.Organization) bool { return o.LegalName != s.systemOrganizationConfig.LegalName },
		} {
			if diff(org) {
				changes = append(changes, check)
				break
			}
		}

		if len(changes) > 0 {
			s.logger.Info("system organization update triggered", slog.Any("changes", changes))

			org.LegalName = s.systemOrganizationConfig.LegalName
			org.UpdateTime = now

			o, uErr := w.Update(ctx, org)

			if uErr != nil {
				return nil, fmt.Errorf("failed to update system organization: %w", uErr)
			}

			return o, nil
		}
	}

	// if the organization does not exist, create it
	if errors.Is(err, storage.ErrOrganizationNotFound) {
		s.logger.Info("system organization creation triggered")

		o, cErr := w.Create(ctx, &domain.Organization{
			ID:         s.systemOrganizationConfig.ID,
			LegalName:  s.systemOrganizationConfig.LegalName,
			Slug:       s.systemOrganizationConfig.Slug,
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
func (s *Initializer) setupUser(ctx context.Context, org *domain.Organization, r *user.Reader, w *user.Writer) (*domain.User, error) {
	currentPassword := []byte(s.systemUserConfig.Password)
	hashedPassword, hErr := password.Hash(currentPassword)
	if hErr != nil {
		return nil, hErr
	}

	// check if the user already exists
	u, err := r.GetByEmail(ctx, s.systemUserConfig.Email)
	if err == nil {
		// check if the user needs to be updated
		changes := []storage.UserField{}

		for check, diff := range map[storage.UserField]func(*domain.User) bool{
			storage.UserOrganizationID: func(u *domain.User) bool { return u.OrganizationID != org.ID },
			storage.UserDisplayName:    func(u *domain.User) bool { return u.DisplayName != systemUserName },
			storage.UserEmail:          func(u *domain.User) bool { return u.Email != s.systemUserConfig.Email },
			storage.UserHashedPassword: func(u *domain.User) bool { return password.Check(currentPassword, []byte(u.HashedPassword)) != nil },
		} {
			if diff(u) {
				changes = append(changes, check)
			}
		}

		if len(changes) > 0 {
			s.logger.Info("system user update triggered", slog.Any("changes", changes))

			u.OrganizationID = org.ID
			u.DisplayName = systemUserName
			u.Email = s.systemUserConfig.Email
			u.HashedPassword = string(hashedPassword)

			u, err = w.Update(ctx, u)
			if err != nil {
				return nil, fmt.Errorf("failed to update system user: %w", err)
			}

			return u, nil
		}
	}

	// if the user does not exist, create it
	if errors.Is(err, storage.ErrUserNotFound) {
		s.logger.Info("system user creation triggered")

		u, err = w.Create(ctx, &domain.User{
			ID:             s.systemUserConfig.ID,
			OrganizationID: org.ID,
			DisplayName:    systemUserName,
			Email:          s.systemUserConfig.Email,
			Status:         domain.UserStatusActive,
			HashedPassword: string(hashedPassword),
		})

		if err != nil {
			return nil, fmt.Errorf("failed to create system user: %w", err)
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
	org, err := s.setupOrganization(
		ctx,
		organization.NewReader(r.Organization),
		organization.NewWriter(s.clock, r.Organization),
	)
	if err != nil {
		return err
	}

	if _, err = s.setupUser(
		ctx,
		org,
		user.NewReader(r.User),
		user.NewWriter(s.clock, r.User),
	); err != nil {
		return err
	}

	return nil
}
