package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/extreme-business/lingo/apps/account/storage"
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
	namePairs := strings.Split(in.GetName(), "/")
	if len(namePairs)%2 != 0 {
		return fmt.Errorf("invalid name %q", in.GetName())
	}

	for i := 0; i < len(namePairs); i += 2 {
		switch namePairs[i] {
		case "organizations":
			organizationID = namePairs[i+1]
		case "users":
			id = namePairs[i+1]
		default:
			return fmt.Errorf("unknown name %q", namePairs[i])
		}
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
	if in.GetDeleteTime() != nil {
		u.DeleteTime = in.GetDeleteTime().AsTime()
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

// ToDomain maps a User to a domain.User.
func (u *User) ToStorage(out *storage.User) error {
	for _, field := range storage.UserFields() {
		switch field {
		case storage.UserID:
			out.ID = u.ID
		case storage.UserOrganizationID:
			out.OrganizationID = u.OrganizationID
		case storage.UserDisplayName:
			out.DisplayName = u.DisplayName
		case storage.UserEmail:
			out.Email = u.Email
		case storage.UserHashedPassword:
			out.HashedPassword = u.HashedPassword
		case storage.UserStatus:
			out.Status = string(u.Status)
		case storage.UserCreateTime:
			out.CreateTime = u.CreateTime
		case storage.UserUpdateTime:
			out.UpdateTime = u.UpdateTime
		case storage.UserDeleteTime:
			out.DeleteTime.Time = u.DeleteTime
			out.DeleteTime.Valid = !u.DeleteTime.IsZero()
		default:
			return fmt.Errorf("unknown field %q", field)
		}
	}

	return nil
}

// FromStorage maps a storage.User to a User.
func (u *User) FromStorage(in *storage.User) error {
	for _, field := range storage.UserFields() {
		switch field {
		case storage.UserID:
			u.ID = in.ID
		case storage.UserOrganizationID:
			u.OrganizationID = in.OrganizationID
		case storage.UserDisplayName:
			u.DisplayName = in.DisplayName
		case storage.UserEmail:
			u.Email = in.Email
		case storage.UserHashedPassword:
			u.HashedPassword = in.HashedPassword
		case storage.UserStatus:
			u.Status = UserStatus(in.Status)
		case storage.UserCreateTime:
			u.CreateTime = in.CreateTime
		case storage.UserUpdateTime:
			u.UpdateTime = in.UpdateTime
		case storage.UserDeleteTime:
			u.DeleteTime = in.DeleteTime.Time
		default:
			return fmt.Errorf("unknown field %q", field)
		}
	}

	return nil
}
