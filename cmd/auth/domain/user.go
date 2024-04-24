package domain

import (
	"fmt"
	"time"

	protoauth "github.com/dwethmar/lingo/proto/gen/go/public/auth/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// User is a user.
type User struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	DisplayName    string
	Email          string
	Password       string
	CreateTime     time.Time
	UpdateTime     time.Time
	Organizations  []*Organization
}

func (u *User) FromProto(in *protoauth.User) error {
	var organizationID string
	var ID string

	if _, err := fmt.Sscanf(in.Name, "organization/%s/user/%s", &organizationID, &ID); err != nil {
		return err
	}

	userUUID, err := uuid.Parse(ID)
	if err != nil {
		return err
	}

	organizationUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return err
	}

	u.ID = userUUID
	u.OrganizationID = organizationUUID
	u.DisplayName = in.DisplayName
	u.Email = in.Email

	if in.CreateTime != nil {
		u.CreateTime = in.CreateTime.AsTime()
	}

	if in.UpdateTime != nil {
		u.UpdateTime = in.UpdateTime.AsTime()
	}

	return nil
}

func (u *User) ToProto(in *protoauth.User) error {
	in.Name = fmt.Sprintf("organization/%s/user/%s", u.OrganizationID.String(), u.ID.String())
	in.DisplayName = u.DisplayName
	in.Email = u.Email
	in.CreateTime = timestamppb.New(u.CreateTime)
	in.UpdateTime = timestamppb.New(u.UpdateTime)
	return nil
}
