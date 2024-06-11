package token_test

import (
	"testing"
	"time"

	"github.com/dwethmar/lingo/cmd/account/token"
)

func TestExpirationTime(t *testing.T) {
	loc, lErr := time.LoadLocation("Europe/Amsterdam")
	if lErr != nil {
		t.Fatal(lErr)
	}

	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "valid token",
			args: args{
				tokenString: string(validToken),
			},
			want:    time.Date(2099, 4, 18, 17, 41, 29, 0, loc),
			wantErr: false,
		},
		{
			name: "invalid token",
			args: args{
				tokenString: "invalid token",
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := token.ExpirationTime(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractExpirationTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !got.Equal(tt.want) {
				t.Errorf("ExtractExpirationTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
