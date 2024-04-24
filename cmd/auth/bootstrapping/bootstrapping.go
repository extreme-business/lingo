// Package bootstrapping contains the Initializer to set up the system data.
package bootstrapping

import (
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/google/uuid"
)

// Initializer is responsible for setting up the system user and organization.
type Initializer struct {
	systemUserID           uuid.UUID // systemUserID is the user id of the system user.
	systemUserEmail        string    // systemUserEmail is the email of the system user.
	systemOrganizationID   uuid.UUID // systemOrganizationID is the organization id of the system user.
	systemOrganizationName string    // systemOrganizationName is the organization name of the system user.
	userRepo               user.Repository
}

type Config struct {
	SystemUserID     uuid.UUID
	SystemUserEmail  string
	OrganizationID   uuid.UUID
	OrganizationName string
	UserRepo         user.Repository
}

func NewInitializer(config Config) *Initializer {
	return &Initializer{
		systemUserID:           config.SystemUserID,
		systemUserEmail:        config.SystemUserEmail,
		systemOrganizationID:   config.OrganizationID,
		systemOrganizationName: config.OrganizationName,
		userRepo:               config.UserRepo,
	}
}

func (s *Initializer) Setup() error {
	return nil
}
