package domain

import (
	"fmt"
	"time"

	protoaccount "github.com/extreme-business/lingo/proto/gen/go/public/account/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserCollection is the name of the user collection.
const UserCollection = "users"

type UserStatus string

func (s UserStatus) String() string { return string(s) }

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusDeleted  UserStatus = "deleted"
)

// User is a user who uses or operates the system.
type User struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	DisplayName    string
	Email          string
	HashedPassword string
	Status         UserStatus
	CreateTime     time.Time
	UpdateTime     time.Time
	DeleteTime     time.Time
	Organization   *Organization // Organization is the primary organization the user belongs to.
}

func (u *User) FromProto(in *protoaccount.User) error {
	var organizationID string
	var id string

	if _, err := fmt.Sscanf(in.GetName(), "organizations/%s/users/%s", &organizationID, &id); err != nil {
		return err
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	organizationUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return err
	}

	u.ID = userUUID
	u.OrganizationID = organizationUUID
	u.DisplayName = in.GetDisplayName()
	u.Email = in.GetEmail()

	if in.GetCreateTime() != nil {
		u.CreateTime = in.GetCreateTime().AsTime()
	}

	if in.GetUpdateTime() != nil {
		u.UpdateTime = in.GetUpdateTime().AsTime()
	}

	return nil
}

func (u *User) ToProto(in *protoaccount.User) error {
	in.Name = fmt.Sprintf("organizations/%s/users/%s", u.OrganizationID.String(), u.ID.String())
	in.DisplayName = u.DisplayName
	in.Email = u.Email
	in.CreateTime = timestamppb.New(u.CreateTime)
	in.UpdateTime = timestamppb.New(u.UpdateTime)
	return nil
}
