package registration

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/dwethmar/lingo/pkg/uuidgen"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNewManager(t *testing.T) {
	t.Run("should create a new manager", func(t *testing.T) {
		c := Config{
			Clock:    clock.Default(),
			UserRepo: nil,
		}
		m := NewManager(c)
		if m == nil {
			t.Fatalf("NewManager() = nil, want a manager")
		}

		expected := &Manager{
			clock:    c.Clock,
			userRepo: c.UserRepo,
		}

		if m.clock != expected.clock {
			t.Fatalf("NewManager() = %v, want %v", m.clock, expected.clock)
		}

		if m.userRepo != expected.userRepo {
			t.Fatalf("NewManager() = %v, want %v", m.userRepo, expected.userRepo)
		}
	})
}

func TestManager_Register(t *testing.T) {
	t.Run("should create a new user", func(t *testing.T) {
		now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		userRepo := user.MockRepository{
			CreateFunc: func(ctx context.Context, u *user.User) (*user.User, error) {
				return &user.User{
					ID:         u.ID,
					Username:   u.Username,
					Email:      u.Email,
					Password:   u.Password,
					CreateTime: u.CreateTime,
					UpdateTime: u.UpdateTime,
				}, nil
			},
		}

		m := Manager{
			uuidgen: uuidgen.New(func() uuid.UUID {
				return uuid.Must(uuid.Parse("c5172a66-3dbe-4415-bbf9-9921d9798698"))
			}),
			clock: clock.New(time.UTC, func() time.Time {
				return now.Add(time.Second)
			}),
			userRepo:              &userRepo,
			registrationValidator: NewRegistrationValidator(),
		}

		u, err := m.Register(context.TODO(), Registration{
			Username: "username",
			Email:    "email",
			Password: "password!1",
		})
		if err != nil {
			t.Fatalf("Register() = %v, want nil", err)
		}

		expected := &domain.User{
			ID:            uuid.Must(uuid.Parse("c5172a66-3dbe-4415-bbf9-9921d9798698")),
			Username:      "username",
			Email:         "email",
			CreateTime:    now.Add(time.Second),
			UpdateTime:    now.Add(time.Second),
			Organizations: nil,
		}

		if u == nil {
			t.Fatalf("Register() = nil, want a user")
		}

		if diff := cmp.Diff(expected, u); diff != "" {
			t.Fatalf("Register() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("should return an error if the registration is invalid", func(t *testing.T) {
		m := Manager{
			registrationValidator: NewRegistrationValidator(),
		}

		_, err := m.Register(context.TODO(), Registration{
			Username: "a", // too short
			Email:    "test@test.com",
			Password: "test-password1",
		})

		if err == nil {
			t.Fatalf("Register() = nil, want an error")
		}
	})

	t.Run("should return an error if the user repo fails", func(t *testing.T) {
		now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		userRepo := user.MockRepository{
			CreateFunc: func(ctx context.Context, u *user.User) (*user.User, error) {
				return nil, errors.New("error")
			},
		}

		m := Manager{
			uuidgen: uuidgen.New(func() uuid.UUID {
				return uuid.Must(uuid.Parse("c5172a66-3dbe-4415-bbf9-9921d9798698"))
			}),
			clock: clock.New(time.UTC, func() time.Time {
				return now.Add(time.Second)
			}),
			userRepo:              &userRepo,
			registrationValidator: NewRegistrationValidator(),
		}

		u, err := m.Register(context.TODO(), Registration{
			Username: "username",
			Email:    "email",
			Password: "password!1",
		})

		if err == nil {
			t.Errorf("Register() = %v, want nil", err)
		}

		if u != nil {
			t.Errorf("Register() = %v, want nil", u)
		}
	})
}
