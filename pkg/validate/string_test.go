package validate

import (
	"testing"
)

func TestStringValidator_Validate(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		v       StringValidator
		args    args
		wantErr bool
	}{
		{
			name:    "should return an error if the string is empty",
			v:       StringValidator{MinLength("test", 1)},
			args:    args{s: ""},
			wantErr: true,
		},
		{
			name:    "should return no error if the string is not empty",
			v:       StringValidator{MinLength("test", 1)},
			args:    args{s: "a"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("StringValidator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	t.Run("should return a StringValidatorFunc that return no error if a string is at least n characters long", func(t *testing.T) {
		v := MinLength("test", 1)
		if err := v("a"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string is not at least n characters long", func(t *testing.T) {
		v := MinLength("test", 1)
		if err := v(""); err == nil {
			t.Error("expected an error")
		}
	})
}

func TestMaxLength(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string is under n characters", func(t *testing.T) {
		v := MaxLength("test", 5)
		if err := v("ab"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string exceeds n characters", func(t *testing.T) {
		v := MaxLength("test", 1)
		if err := v("ab"); err == nil {
			t.Error("expected an error")
		}
	})
}

func TestContainsSpecialChars(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string contains n special characters", func(t *testing.T) {
		v := ContainsSpecialChars("test", 1)
		if err := v("a!b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		v = ContainsSpecialChars("test", 2)
		if err := v("a!!b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string does not contain n special characters", func(t *testing.T) {
		v := ContainsSpecialChars("test", 1)
		if err := v("ab"); err == nil {
			t.Error("expected an error")
		}
	})
}

func TestContainsDigits(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string contains n digits", func(t *testing.T) {
		v := ContainsDigits("test", 1)
		if err := v("a1b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		v = ContainsDigits("test", 2)
		if err := v("a12b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string does not contain n digits", func(t *testing.T) {
		v := ContainsDigits("test", 1)
		if err := v("ab"); err == nil {
			t.Error("expected an error")
		}
	})

}

func TestSpecialCharWhitelist(t *testing.T) {
	t.Run("should return a StringValidatorFunc that returns no error if a string contains only special characters in the whitelist", func(t *testing.T) {
		v := SpecialCharWhitelist("test", 1, '!')
		if err := v("a!b"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should return a StringValidatorFunc that returns an error if a string contains a special character not in the whitelist", func(t *testing.T) {
		v := SpecialCharWhitelist("test", 1, '!')
		if err := v("a@b"); err == nil {
			t.Error("expected an error")
		}
	})
}
