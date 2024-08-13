package config_test

import (
	"testing"

	"github.com/extreme-business/lingo/pkg/config"
	"github.com/spf13/viper"
)

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		got := config.New()
		if got == nil {
			t.Errorf("New() = %v, want %v", got, nil)
		}
	})
}

func TestConfig_DatabaseURL(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_DB_URL is not set", func(t *testing.T) {
		v, err := config.New().DatabaseURL()
		if err == nil {
			t.Errorf("DatabaseURL() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("DatabaseURL() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_DB_URL", func(t *testing.T) {
		t.Setenv("LINGO_DB_URL", "test")
		c := config.New()
		got, err := c.DatabaseURL()
		if err != nil {
			t.Errorf("DatabaseURL() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("DatabaseURL() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_SigningKeyRegistration(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_SIGNING_KEY_REGISTRATION is not set", func(t *testing.T) {
		v, err := config.New().SigningKeyRegistration()
		if err == nil {
			t.Errorf("SigningKeyRegistration() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("SigningKeyRegistration() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_SIGNING_KEY_REGISTRATION", func(t *testing.T) {
		t.Setenv("LINGO_SIGNING_KEY_REGISTRATION", "test")
		c := config.New()
		got, err := c.SigningKeyRegistration()
		if err != nil {
			t.Errorf("SigningKeyRegistration() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("SigningKeyRegistration() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_SigningKeyAuthentication(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_SIGNING_KEY_AUTHENTICATION is not set", func(t *testing.T) {
		v, err := config.New().SigningKeyAuthentication()
		if err == nil {
			t.Errorf("SigningKeyAuthentication() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("SigningKeyAuthentication() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_SIGNING_KEY_AUTHENTICATION", func(t *testing.T) {
		t.Setenv("LINGO_SIGNING_KEY_AUTHENTICATION", "test")
		c := config.New()
		got, err := c.SigningKeyAuthentication()
		if err != nil {
			t.Errorf("SigningKeyAuthentication() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("SigningKeyAuthentication() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_HTTPPort(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_HTTP_PORT is not set", func(t *testing.T) {
		v, err := config.New().HTTPPort()
		if err == nil {
			t.Errorf("HTTPPort() error = %v, wantErr %v", err, true)
		}

		if v != 0 {
			t.Errorf("HTTPPort() = %v, want %v", v, 0)
		}
	})

	t.Run("should return the value of LINGO_HTTP_PORT", func(t *testing.T) {
		t.Setenv("LINGO_HTTP_PORT", "8080")
		c := config.New()
		got, err := c.HTTPPort()
		if err != nil {
			t.Errorf("HTTPPort() error = %v, wantErr %v", err, nil)
		}
		if got != 8080 {
			t.Errorf("HTTPPort() = %v, want %v", got, 8080)
		}
	})
}

func TestConfig_GRPCPort(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_GRPC_PORT is not set", func(t *testing.T) {
		v, err := config.New().GRPCPort()
		if err == nil {
			t.Errorf("GRPCPort() error = %v, wantErr %v", err, true)
		}

		if v != 0 {
			t.Errorf("GRPCPort() = %v, want %v", v, 0)
		}
	})

	t.Run("should return the value of LINGO_GRPC_PORT", func(t *testing.T) {
		t.Setenv("LINGO_GRPC_PORT", "8080")
		c := config.New()
		got, err := c.GRPCPort()
		if err != nil {
			t.Errorf("GRPCPort() error = %v, wantErr %v", err, nil)
		}
		if got != 8080 {
			t.Errorf("GRPCPort() = %v, want %v", got, 8080)
		}
	})
}

func TestConfig_GrpcTLSCertFile(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_GRPC_TLS_CERT_FILE is not set", func(t *testing.T) {
		v, err := config.New().GrpcTLSCertFile()
		if err == nil {
			t.Errorf("GrpcTLSCertFile() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("GrpcTLSCertFile() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_GRPC_TLS_CERT_FILE", func(t *testing.T) {
		t.Setenv("LINGO_GRPC_TLS_CERT_FILE", "test")
		c := config.New()
		got, err := c.GrpcTLSCertFile()
		if err != nil {
			t.Errorf("GrpcTLSCertFile() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("GrpcTLSCertFile() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_GrpcTLSKeyFile(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_GRPC_TLS_KEY_FILE is not set", func(t *testing.T) {
		v, err := config.New().GrpcTLSKeyFile()
		if err == nil {
			t.Errorf("GrpcTLSKeyFile() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("GrpcTLSKeyFile() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_GRPC_TLS_KEY_FILE", func(t *testing.T) {
		t.Setenv("LINGO_GRPC_TLS_KEY_FILE", "test")
		c := config.New()
		got, err := c.GrpcTLSKeyFile()
		if err != nil {
			t.Errorf("GrpcTLSKeyFile() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("GrpcTLSKeyFile() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_HTTPTLSKeyFile(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_HTTP_TLS_KEY_FILE is not set", func(t *testing.T) {
		v, err := config.New().HTTPTLSKeyFile()
		if err == nil {
			t.Errorf("HTTPTLSKeyFile() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("HTTPTLSKeyFile() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_HTTP_TLS_KEY_FILE", func(t *testing.T) {
		t.Setenv("LINGO_HTTP_TLS_KEY_FILE", "test")
		c := config.New()
		got, err := c.HTTPTLSKeyFile()
		if err != nil {
			t.Errorf("HTTPTLSKeyFile() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("HTTPTLSKeyFile() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_HTTPTLSCertFile(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_HTTP_TLS_CERT_FILE is not set", func(t *testing.T) {
		v, err := config.New().HTTPTLSCertFile()
		if err == nil {
			t.Errorf("HTTPTLSCertFile() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("HTTPTLSCertFile() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_HTTP_TLS_CERT_FILE", func(t *testing.T) {
		t.Setenv("LINGO_HTTP_TLS_CERT_FILE", "test")
		c := config.New()
		got, err := c.HTTPTLSCertFile()
		if err != nil {
			t.Errorf("HTTPTLSCertFile() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("HTTPTLSCertFile() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_AccountURL(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_AUTH_URL is not set", func(t *testing.T) {
		v, err := config.New().AccountServiceURL()
		if err == nil {
			t.Errorf("AccountUrl() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("AccountUrl() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_ACCOUNT_SERVICE_URL", func(t *testing.T) {
		t.Setenv("LINGO_ACCOUNT_SERVICE_URL", "test")
		c := config.New()
		got, err := c.AccountServiceURL()
		if err != nil {
			t.Errorf("AccountUrl() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("AccountUrl() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_SystemUserID(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_SYSTEM_USER_ID is not set", func(t *testing.T) {
		v, err := config.New().SystemUserID()
		if err == nil {
			t.Errorf("SystemUserID() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("SystemUserID() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_SYSTEM_USER_ID", func(t *testing.T) {
		t.Setenv("LINGO_SYSTEM_USER_ID", "test")
		c := config.New()
		got, err := c.SystemUserID()
		if err != nil {
			t.Errorf("SystemUserID() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("SystemUserID() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_SystemUserEmail(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_SYSTEM_USER_EMAIL is not set", func(t *testing.T) {
		v, err := config.New().SystemUserEmail()
		if err == nil {
			t.Errorf("SystemUserEmail() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("SystemUserEmail() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_SYSTEM_USER_EMAIL", func(t *testing.T) {
		t.Setenv("LINGO_SYSTEM_USER_EMAIL", "test")
		c := config.New()
		got, err := c.SystemUserEmail()
		if err != nil {
			t.Errorf("SystemUserEmail() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("SystemUserEmail() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_SystemUserPassword(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_SYSTEM_USER_PASSWORD is not set", func(t *testing.T) {
		v, err := config.New().SystemUserPassword()
		if err == nil {
			t.Errorf("SystemUserPassword() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("SystemUserPassword() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_SYSTEM_USER_PASSWORD", func(t *testing.T) {
		t.Setenv("LINGO_SYSTEM_USER_PASSWORD", "test")
		c := config.New()
		got, err := c.SystemUserPassword()
		if err != nil {
			t.Errorf("SystemUserPassword() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("SystemUserPassword() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_SystemOrganizationID(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_SYSTEM_ORGANIZATION_ID is not set", func(t *testing.T) {
		v, err := config.New().SystemOrganizationID()
		if err == nil {
			t.Errorf("SystemOrganizationID() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("SystemOrganizationID() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_SYSTEM_ORGANIZATION_ID", func(t *testing.T) {
		t.Setenv("LINGO_SYSTEM_ORGANIZATION_ID", "test")
		c := config.New()
		got, err := c.SystemOrganizationID()
		if err != nil {
			t.Errorf("SystemOrganizationID() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("SystemOrganizationID() = %v, want %v", got, "test")
		}
	})
}

func TestConfig_SystemOrganizationLegalName(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	t.Run("should return error if LINGO_SYSTEM_ORGANIZATION_LEGAL_NAME is not set", func(t *testing.T) {
		v, err := config.New().SystemOrganizationLegalName()
		if err == nil {
			t.Errorf("SystemOrganizationLegalName() error = %v, wantErr %v", err, true)
		}

		if v != "" {
			t.Errorf("SystemOrganizationLegalName() = %v, want %v", v, "")
		}
	})

	t.Run("should return the value of LINGO_SYSTEM_ORGANIZATION_LEGAL_NAME", func(t *testing.T) {
		t.Setenv("LINGO_SYSTEM_ORGANIZATION_LEGAL_NAME", "test")
		c := config.New()
		got, err := c.SystemOrganizationLegalName()
		if err != nil {
			t.Errorf("SystemOrganizationLegalName() error = %v, wantErr %v", err, nil)
		}
		if got != "test" {
			t.Errorf("SystemOrganizationLegalName() = %v, want %v", got, "test")
		}
	})
}
