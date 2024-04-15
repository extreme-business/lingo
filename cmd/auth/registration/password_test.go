package registration

import "testing"

func TestHashPassword(t *testing.T) {
	t.Run("should hash a password", func(t *testing.T) {
		password := "password"
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword() error = %v", err)
		}

		if !CheckPasswordHash(password, hash) {
			t.Fatalf("CheckPasswordHash() = false, want true")
		}
	})
}

func TestValidatePassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return an error if the password is too short",
			args: args{
				password: "1234567",
			},
			wantErr: true,
		},
		{
			name: "should not return an error if the password is long enough",
			args: args{
				password: "12345678",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
