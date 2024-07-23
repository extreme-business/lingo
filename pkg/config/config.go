package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	// env keys.
	envKeyDatabaseURL              = "DB_URL"
	envKeySigningKeyRegistration   = "SIGNING_KEY_REGISTRATION"
	envKeySigningKeyAuthentication = "SIGNING_KEY_AUTHENTICATION"
	envKeyHTTPPort                 = "HTTP_PORT"
	envKeyGRPCPort                 = "GRPC_PORT"
	envKeyGrpcTLSCertFile          = "GRPC_TLS_CERT_FILE"
	envKeyGrpcTLSKeyFile           = "GRPC_TLS_KEY_FILE"
	envKeyHTTPTLSKeyFile           = "HTTP_TLS_KEY_FILE"
	envKeyHTTPTLSCertFile          = "HTTP_TLS_CERT_FILE"
	envKeyAccountServiceURL        = "ACCOUNT_SERVICE_URL"
	envSystemUserID                = "SYSTEM_USER_ID"
	envSystemUserEmail             = "SYSTEM_USER_EMAIL"
	envSystemUserPassword          = "SYSTEM_USER_PASSWORD"
	envSystemOrganizationID        = "SYSTEM_ORGANIZATION_ID"
	envSystemOrganizationLegalName = "SYSTEM_ORGANIZATION_LEGAL_NAME"
	envSystemOrganizationSlug      = "SYSTEM_ORGANIZATION_SLUG"
)

type Config struct {
	viper *viper.Viper
}

func New() *Config {
	v := viper.New()
	v.SetEnvPrefix("lingo")
	v.AutomaticEnv()

	return &Config{
		viper: v,
	}
}

// getString returns the value of the key as a string.
func (c *Config) getString(key string) (string, error) {
	if !c.viper.IsSet(key) {
		return "", fmt.Errorf("%s is not set", key)
	}

	return c.viper.GetString(key), nil
}

// getInt returns the value of the key as an integer.
func (c *Config) getInt(key string) (int, error) {
	if !c.viper.IsSet(key) {
		return 0, fmt.Errorf("%s is not set", key)
	}

	return c.viper.GetInt(key), nil
}

func (c *Config) DatabaseURL() (string, error) { return c.getString(envKeyDatabaseURL) }
func (c *Config) SigningKeyRegistration() (string, error) {
	return c.getString(envKeySigningKeyRegistration)
}
func (c *Config) SigningKeyAuthentication() (string, error) {
	return c.getString(envKeySigningKeyAuthentication)
}
func (c *Config) HTTPPort() (int, error)                { return c.getInt(envKeyHTTPPort) }
func (c *Config) GRPCPort() (int, error)                { return c.getInt(envKeyGRPCPort) }
func (c *Config) GrpcTLSCertFile() (string, error)      { return c.getString(envKeyGrpcTLSCertFile) }
func (c *Config) GrpcTLSKeyFile() (string, error)       { return c.getString(envKeyGrpcTLSKeyFile) }
func (c *Config) HTTPTLSKeyFile() (string, error)       { return c.getString(envKeyHTTPTLSKeyFile) }
func (c *Config) HTTPTLSCertFile() (string, error)      { return c.getString(envKeyHTTPTLSCertFile) }
func (c *Config) AccountServiceURL() (string, error)    { return c.getString(envKeyAccountServiceURL) }
func (c *Config) SystemUserID() (string, error)         { return c.getString(envSystemUserID) }
func (c *Config) SystemUserEmail() (string, error)      { return c.getString(envSystemUserEmail) }
func (c *Config) SystemUserPassword() (string, error)   { return c.getString(envSystemUserPassword) }
func (c *Config) SystemOrganizationID() (string, error) { return c.getString(envSystemOrganizationID) }
func (c *Config) SystemOrganizationLegalName() (string, error) {
	return c.getString(envSystemOrganizationLegalName)
}
func (c *Config) SystemOrganizationSlug() (string, error) {
	return c.getString(envSystemOrganizationSlug)
}
