package user_test

import (
	"testing"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestUser_ToDomain(t *testing.T) {
	t.Run("should map a user to a domain user", func(t *testing.T) {
		u := &user.User{
			ID:          uuid.Must(uuid.Parse("c5172a66-3dbe-4415-bbf9-9921d9798698")),
			DisplayName: "username",
			Email:       "email",
			Password:    "password",
			CreateTime:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		d := &domain.User{}
		u.ToDomain(d)

		expected := &domain.User{
			ID:          uuid.Must(uuid.Parse("c5172a66-3dbe-4415-bbf9-9921d9798698")),
			DisplayName: "username",
			Email:       "email",
			Password:    "password",
			CreateTime:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if diff := cmp.Diff(expected, d); diff != "" {
			t.Errorf("ToDomain() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestUser_FromDomain(t *testing.T) {
	t.Run("should map a domain user to a user", func(t *testing.T) {
		d := &domain.User{
			ID:          uuid.Must(uuid.Parse("c5172a66-3dbe-4415-bbf9-9921d9798698")),
			DisplayName: "username",
			Email:       "email",
			Password:    "password",
			CreateTime:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		u := &user.User{}
		u.FromDomain(d)

		expected := &user.User{
			ID:          uuid.Must(uuid.Parse("c5172a66-3dbe-4415-bbf9-9921d9798698")),
			DisplayName: "username",
			Email:       "email",
			Password:    "password",
			CreateTime:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if diff := cmp.Diff(expected, u); diff != "" {
			t.Errorf("FromDomain() mismatch (-want +got):\n%s", diff)
		}
	})
}
