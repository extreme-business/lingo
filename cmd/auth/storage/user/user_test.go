package user

import (
	"testing"
	"time"

	"github.com/dwethmar/lingo/cmd/auth/domain"
	"github.com/google/uuid"
)

func TestUser_ToDomain(t *testing.T) {
	t.Run("should map a user to a domain user", func(t *testing.T) {})
}

func TestUser_FromDomain(t *testing.T) {
	type fields struct {
		ID         uuid.UUID
		Username   string
		Email      string
		Password   string
		CreateTime time.Time
		UpdateTime time.Time
	}
	type args struct {
		in *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:         tt.fields.ID,
				Username:   tt.fields.Username,
				Email:      tt.fields.Email,
				Password:   tt.fields.Password,
				CreateTime: tt.fields.CreateTime,
				UpdateTime: tt.fields.UpdateTime,
			}
			if err := u.FromDomain(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("User.FromDomain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
