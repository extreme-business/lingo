package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/extreme-business/lingo/cmd/account/domain"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
	// Update.
	ErrNoUserFieldsToUpdate = errors.New("no fields to update")
	// Fields.
	ErrUserUnknownField = errors.New("unknown user field")
	// sort errors.
	ErrEmptyUserSortField       = errors.New("user field is empty")
	ErrInvalidUserSortDirection = errors.New("invalid user sort direction")
	// Unique constraint errors.
	ErrConflictUserID          = errors.New("unique id conflict")
	ErrConflictUserDisplayName = errors.New("unique display_name conflict")
	ErrConflictUserEmail       = errors.New("unique email conflict")
	// Immutable errors.
	ErrImmutableUserID         = errors.New("field id is read-only")
	ErrImmutableUserCreateTime = errors.New("field create_time is read-only")
)

type UserField string

const (
	UserID             UserField = "id"
	UserOrganizationID UserField = "organization_id"
	UserDisplayName    UserField = "display_name"
	UserEmail          UserField = "email"
	UserPassword       UserField = "password"
	UserCreateTime     UserField = "create_time"
	UserUpdateTime     UserField = "update_time"
)

// UserFields returns all user fields.
func UserFields() []UserField {
	return []UserField{
		UserID,
		UserOrganizationID,
		UserDisplayName,
		UserEmail,
		UserPassword,
		UserCreateTime,
		UserUpdateTime,
	}
}

type User struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	DisplayName    string
	Email          string
	Password       string
	CreateTime     time.Time
	UpdateTime     time.Time
}

// ToDomain maps a User to a domain.User.
func (u *User) ToDomain(in *domain.User) {
	in.ID = u.ID
	in.OrganizationID = u.OrganizationID
	in.DisplayName = u.DisplayName
	in.Email = u.Email
	in.Password = u.Password
	in.CreateTime = u.CreateTime
	in.UpdateTime = u.UpdateTime
}

// FromDomain maps a domain.User to a User.
func (u *User) FromDomain(in *domain.User) {
	u.ID = in.ID
	u.OrganizationID = in.OrganizationID
	u.DisplayName = in.DisplayName
	u.Email = in.Email
	u.Password = in.Password
	u.CreateTime = in.CreateTime
	u.UpdateTime = in.UpdateTime
}

// UserSort pairs a field with a direction.
type UserSort struct {
	Field     UserField
	Direction Direction
}

type UserOrderBy []UserSort

// Validate checks if the sort fields are valid.
func (o UserOrderBy) Validate() error {
	fields := UserFields()
	for _, s := range o {
		if s.Field == "" {
			return ErrEmptyUserSortField
		}

		var found bool
		for _, f := range fields {
			if s.Field == f {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("%s: %w", s.Field, ErrUserUnknownField)
		}

		if s.Direction != ASC && s.Direction != DESC {
			return fmt.Errorf("%s: %w", s.Direction, ErrInvalidUserSortDirection)
		}
	}

	return nil
}

type UserReader interface {
	Get(context.Context, uuid.UUID) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	List(context.Context, Pagination, UserOrderBy, ...Condition) ([]*User, error)
}

type UserWriter interface {
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User, []UserField) (*User, error)
	Delete(context.Context, uuid.UUID) error
}

// UserRepository is a reader and writer for users.
type UserRepository interface {
	UserReader
	UserWriter
}

// UserByOrganizationIDCondition is a search condition for users by organization ID.
type UserByOrganizationIDCondition struct {
	OrganizationID uuid.UUID
}

func (UserByOrganizationIDCondition) condition() {}
