package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserError error

var (
	ErrUserNotFound UserError = errors.New("user not found")
	// Update.
	ErrNoUserFieldsToUpdate UserError = errors.New("no fields to update")
	// Fields.
	ErrUserUnknownField UserError = errors.New("unknown user field")
	// sort errors.
	ErrEmptyUserSortField       UserError = errors.New("user field is empty")
	ErrInvalidUserSortDirection UserError = errors.New("invalid user sort direction")
	// Unique constraint errors.
	ErrConflictUserID          UserError = errors.New("unique id conflict")
	ErrConflictUserDisplayName UserError = errors.New("unique display_name conflict")
	ErrConflictUserEmail       UserError = errors.New("unique email conflict")
	// Immutable errors.
	ErrImmutableUserID         UserError = errors.New("field id is read-only")
	ErrImmutableUserCreateTime UserError = errors.New("field create_time is read-only")
)

type UserField string

const (
	UserID             UserField = "id"
	UserOrganizationID UserField = "organization_id"
	UserDisplayName    UserField = "display_name"
	UserEmail          UserField = "email"
	UserHashedPassword UserField = "hashed_password"
	UserStatus         UserField = "status"
	UserCreateTime     UserField = "create_time"
	UserUpdateTime     UserField = "update_time"
	UserDeleteTime     UserField = "delete_time"
)

// UserFields returns all user fields.
func UserFields() []UserField {
	return []UserField{
		UserID,
		UserOrganizationID,
		UserDisplayName,
		UserEmail,
		UserHashedPassword,
		UserStatus,
		UserCreateTime,
		UserUpdateTime,
		UserDeleteTime,
	}
}

type User struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	DisplayName    string
	Email          string
	HashedPassword string
	Status         string
	CreateTime     time.Time
	UpdateTime     time.Time
	DeleteTime     sql.NullTime
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
