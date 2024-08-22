// Package bootstrapping contains the Initializer to set up the system data.
package bootstrapping

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
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

// Bootstrapper is responsible for setting up the system user and organization.
type Bootstrapper struct {
	logger    *slog.Logger
	clock     func() time.Time
	dbManager storage.DBManager
}

type Config struct {
	Logger    *slog.Logger
	Clock     func() time.Time
	DBManager storage.DBManager
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

	return nil
}

func New(config Config) (*Bootstrapper, error) {
	return &Bootstrapper{
		logger:    config.Logger,
		clock:     config.Clock,
		dbManager: config.DBManager,
	}, config.Validate()
}

// Setup sets up the system user and organization.
func (s *Bootstrapper) Setup(ctx context.Context, systemUserConfig SystemUserConfig, systemOrganizationConfig SystemOrgConfig) error {
	return s.dbManager.BeginOp(ctx, func(ctx context.Context, r storage.Repositories) error {
		if r.User == nil {
			return errors.New("user repository is required")
		}

		if r.Organization == nil {
			return errors.New("organization repository is required")
		}

		return s.setup(ctx, &r, systemUserConfig, systemOrganizationConfig)
	})
}

// setupOrganization sets up the system organization. If the organization already exists, it will be updated if necessary.
func (s *Bootstrapper) setupOrganization(ctx context.Context, r *organization.Reader, w *organization.Writer, c SystemOrgConfig) (*domain.Organization, error) {
	// check if the organization already exists
	org, err := r.Get(ctx, c.ID)
	if err == nil {
		// check if the organization needs to be updated
		changes := []storage.OrganizationField{}
		type check struct {
			field storage.OrganizationField
			diff  func(*domain.Organization) bool
		}
		for _, check := range []check{
			{storage.OrganizationLegalName, func(o *domain.Organization) bool { return o.LegalName != c.LegalName }},
			{storage.OrganizationSlug, func(o *domain.Organization) bool { return o.Slug != c.Slug }},
			{storage.OrganizationUpdateTime, func(o *domain.Organization) bool { return false }},
			{storage.OrganizationCreateTime, func(o *domain.Organization) bool { return false }},
		} {
			if check.diff(org) {
				changes = append(changes, check.field)
			}
		}

		if len(changes) > 0 {
			s.logger.Info("system organization update triggered", slog.Any("changes", changes))
			org.LegalName = c.LegalName
			org.Slug = c.Slug

			o, uErr := w.Update(ctx, org, changes)
			if uErr != nil {
				return nil, fmt.Errorf("failed to update system organization: %w", uErr)
			}

			return o, nil
		}
	}

	// if the organization does not exist, create it
	if errors.Is(err, organization.ErrOrganizationNotFound) {
		s.logger.Info("system organization creation triggered")
		now := s.clock()

		o, cErr := w.Create(ctx, &domain.Organization{
			ID:         c.ID,
			LegalName:  c.LegalName,
			Slug:       c.Slug,
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
func (s *Bootstrapper) setupUser(ctx context.Context, org *domain.Organization, r *user.Reader, w *user.Writer, c SystemUserConfig) (*domain.User, error) {
	currentPassword := []byte(c.Password)
	hashedPassword, hErr := password.Hash(currentPassword)
	if hErr != nil {
		return nil, hErr
	}

	u, err := r.Get(ctx, c.ID)
	if err == nil {
		if u.ID != c.ID {
			return nil, fmt.Errorf("system user id mismatch: expected %s, got %s", c.ID, u.ID)
		}
		// check if the user needs to be updated
		changes := []storage.UserField{}
		type check struct {
			field storage.UserField
			diff  func(*domain.User) bool
		}
		for _, check := range []check{
			{storage.UserOrganizationID, func(u *domain.User) bool { return u.OrganizationID != org.ID }},
			{storage.UserDisplayName, func(u *domain.User) bool { return u.DisplayName != systemUserName }},
			{storage.UserEmail, func(u *domain.User) bool { return u.Email != c.Email }},
			{storage.UserHashedPassword, func(u *domain.User) bool { return password.Check(currentPassword, []byte(u.HashedPassword)) != nil }},
			{storage.UserUpdateTime, func(u *domain.User) bool { return false }},
			{storage.UserCreateTime, func(u *domain.User) bool { return false }},
		} {
			if check.diff(u) {
				changes = append(changes, check.field)
			}
		}

		if len(changes) > 0 {
			slices.Sort(changes)
			s.logger.Info("system user update triggered", slog.Any("changes", changes))
			u.OrganizationID = org.ID
			u.DisplayName = systemUserName
			u.Email = c.Email
			u.HashedPassword = string(hashedPassword)

			u, err = w.Update(ctx, u, changes)
			if err != nil {
				return nil, fmt.Errorf("failed to update system user: %w", err)
			}

			return u, nil
		}
	}

	// if the user does not exist, create it
	if errors.Is(err, user.ErrUserNotFound) {
		s.logger.Info("system user creation triggered")
		now := s.clock()
		u, err = w.Create(ctx, &domain.User{
			ID:             c.ID,
			OrganizationID: org.ID,
			DisplayName:    systemUserName,
			Email:          c.Email,
			Status:         domain.UserStatusActive,
			HashedPassword: string(hashedPassword),
			CreateTime:     now,
			UpdateTime:     now,
			DeleteTime:     time.Time{},
		})

		if err != nil {
			return nil, fmt.Errorf("failed to create system user: %w", err)
		}

		return u, nil
	}

	return nil, err
}

func (s *Bootstrapper) setup(ctx context.Context, r *storage.Repositories, u SystemUserConfig, o SystemOrgConfig) error {
	if s.clock == nil {
		return errors.New("clock is required")
	}

	if err := u.Validate(); err != nil {
		return fmt.Errorf("invalid system user config: %w", err)
	}

	if err := o.Validate(); err != nil {
		return fmt.Errorf("invalid system organization config: %w", err)
	}

	// Create the system organization and user.
	org, err := s.setupOrganization(
		ctx,
		organization.NewReader(r.Organization),
		organization.NewWriter(s.clock, r.Organization),
		o,
	)
	if err != nil {
		return err
	}

	if _, err = s.setupUser(
		ctx,
		org,
		user.NewReader(r.User),
		user.NewWriter(s.clock, r.User),
		u,
	); err != nil {
		return err
	}

	return nil
}
