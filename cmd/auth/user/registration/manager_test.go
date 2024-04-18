package registration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/dwethmar/lingo/cmd/auth/storage/user"
	"github.com/dwethmar/lingo/cmd/auth/user/registration"
	"github.com/dwethmar/lingo/pkg/clock"
	"github.com/dwethmar/lingo/pkg/uuidgen"
	"github.com/dwethmar/lingo/pkg/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNewManager(t *testing.T) {
	t.Run("should create a new manager", func(t *testing.T) {
		c := registration.Config{}
		m := registration.NewManager(c)
		if m == nil {
			t.Fatalf("NewManager() = nil, want a manager")
		}
	})
}

func TestManager_Register(t *testing.T) {
	t.Run("should create a new user", func(t *testing.T) {
		now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		userRepo := user.MockRepository{
			CreateFunc: func(_ context.Context, u *user.User) (*user.User, error) {
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

		m := registration.NewManager(registration.Config{
			UserRepo: &userRepo,
			Clock:    clock.New(time.UTC, func() time.Time { return now.Add(time.Second) }),
			UUIDgen: uuidgen.New(func() uuid.UUID {
				return uuid.Must(uuid.Parse("c5172a66-3dbe-4415-bbf9-9921d9798698"))
			}),
		})

		u, err := m.Register(context.TODO(), registration.Registration{
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

	t.Run("should return an error if registration is invalid", func(t *testing.T) {
		type fields struct {
			config registration.Config
		}
		type args struct {
			ctx          context.Context
			registration registration.Registration
		}
		tests := []struct {
			name   string
			fields fields
			args   args
			want   *domain.User
			want2  string
		}{
			{
				name: "username too short",
				args: args{
					ctx: context.TODO(),
					registration: registration.Registration{
						Username: "a",
						Email:    "email",
						Password: "password!1",
					},
				},
				want:  nil,
				want2: "username",
			},
			{
				name: "username too long",
				args: args{
					ctx: context.TODO(),
					registration: registration.Registration{
						Username: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz",
						Email:    "email",
						Password: "password!1",
					},
				},
				want:  nil,
				want2: "username",
			},
			{
				name: "username contains non allowed special char",
				args: args{
					ctx: context.TODO(),
					registration: registration.Registration{
						Username: "username!",
						Email:    "email",
						Password: "password!1",
					},
				},
				want:  nil,
				want2: "username",
			},
			{
				name: "email too short",
				args: args{
					ctx: context.TODO(),
					registration: registration.Registration{
						Username: "abcdefghijklmnopqrstuvwxyz",
						Email:    "a",
						Password: "password!1",
					},
				},
				want:  nil,
				want2: "email",
			},
			{
				name: "email too long",
				args: args{
					ctx: context.TODO(),
					registration: registration.Registration{
						Username: "abcdefghijklmnopqrstuvwxyz",
						Email:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz",
						Password: "password!1",
					},
				},
				want:  nil,
				want2: "email",
			},
			{
				name: "password too short",
				args: args{
					ctx: context.TODO(),
					registration: registration.Registration{
						Username: "abcdefghijklmnopqrstuvwxyz",
						Email:    "email",
						Password: "a",
					},
				},
				want:  nil,
				want2: "password",
			},
			{
				name: "password too long",
				args: args{
					ctx: context.TODO(),
					registration: registration.Registration{
						Username: "abcdefghijklmnopqrstuvwxyz",
						Email:    "email",
						Password: "1@Aabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz",
					},
				},
				want:  nil,
				want2: "password",
			},
			{
				name: "password does not contain special char",
				args: args{
					ctx: context.TODO(),
					registration: registration.Registration{
						Username: "abcdefghijklmnopqrstuvwxyz",
						Email:    "email",
						Password: "password1",
					},
				},
				want:  nil,
				want2: "password",
			},
			{
				name: "password does not contain digit",
				args: args{
					ctx: context.TODO(),
					registration: registration.Registration{
						Username: "abcdefghijklmnopqrstuvwxyz",
						Email:    "email",
						Password: "password!",
					},
				},
				want:  nil,
				want2: "password",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				m := registration.NewManager(tt.fields.config)
				user, err := m.Register(tt.args.ctx, tt.args.registration)

				if diff := cmp.Diff(tt.want, user); diff != "" {
					t.Fatalf("Register() mismatch (-want +got):\n%s", diff)
				}

				var vErr *validate.Error
				if !errors.As(err, &vErr) {
					t.Fatalf("Register() = %v, want a validate.Error", err)
				}

				if vErr.Field() != tt.want2 {
					t.Errorf("Register() field = %v, want %s", vErr.Field(), tt.want2)
				}
			})
		}
	})

	t.Run("should return an error if the user repo fails", func(t *testing.T) {
		now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		userRepo := user.MockRepository{
			CreateFunc: func(_ context.Context, _ *user.User) (*user.User, error) {
				return nil, errors.New("error")
			},
		}

		m := registration.NewManager(registration.Config{
			UserRepo: &userRepo,
			Clock:    clock.New(time.UTC, func() time.Time { return now }),
			UUIDgen: uuidgen.New(func() uuid.UUID {
				return uuid.Must(uuid.Parse("c5172a66-3dbe-4415-bbf9-9921d9798698"))
			}),
		})

		u, err := m.Register(context.TODO(), registration.Registration{
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
