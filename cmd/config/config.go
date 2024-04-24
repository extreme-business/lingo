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
	envKeyAuthURL                  = "AUTH_URL"
	envSystemUserID                = "SYSTEM_USER_ID"
	envSystemUserEmail             = "SYSTEM_USER_EMAIL"
	envSystemUserPassword          = "SYSTEM_USER_PASSWORD"
	envSystemOrganizationID        = "SYSTEM_ORGANIZATION_ID"
	envSystemOrganizationName      = "SYSTEM_ORGANIZATION_NAME"
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

// configString returns the value of the key as a string.
func (c *Config) configString(key string) (string, error) {
	if !c.viper.IsSet(key) {
		return "", fmt.Errorf("%s is not set", key)
	}

	return c.viper.GetString(key), nil
}

// configInt returns the value of the key as an integer.
func (c *Config) configInt(key string) (int, error) {
	if !c.viper.IsSet(key) {
		return 0, fmt.Errorf("%s is not set", key)
	}

	return c.viper.GetInt(key), nil
}

func (c *Config) DatabaseURL() (string, error) {
	return c.configString(envKeyDatabaseURL)
}

func (c *Config) SigningKeyRegistration() (string, error) {
	return c.configString(envKeySigningKeyRegistration)
}

func (c *Config) SigningKeyAuthentication() (string, error) {
	return c.configString(envKeySigningKeyAuthentication)
}

func (c *Config) HTTPPort() (int, error) {
	return c.configInt(envKeyHTTPPort)
}

func (c *Config) GRPCPort() (int, error) {
	return c.configInt(envKeyGRPCPort)
}

func (c *Config) GrpcTLSCertFile() (string, error) {
	return c.configString(envKeyGrpcTLSCertFile)
}

func (c *Config) GrpcTLSKeyFile() (string, error) {
	return c.configString(envKeyGrpcTLSKeyFile)
}

func (c *Config) HTTPTLSKeyFile() (string, error) {
	return c.configString(envKeyHTTPTLSKeyFile)
}

func (c *Config) HTTPTLSCertFile() (string, error) {
	return c.configString(envKeyHTTPTLSCertFile)
}

func (c *Config) AuthURL() (string, error) {
	return c.configString(envKeyAuthURL)
}

func (c *Config) SystemUserID() (string, error) {
	return c.configString(envSystemUserID)
}

func (c *Config) SystemUserEmail() (string, error) {
	return c.configString(envSystemUserEmail)
}

func (c *Config) SystemUserPassword() (string, error) {
	return c.configString(envSystemUserPassword)
}

func (c *Config) SystemOrganizationID() (string, error) {
	return c.configString(envSystemOrganizationID)
}

func (c *Config) SystemOrganizationName() (string, error) {
	return c.configString(envSystemOrganizationName)
}
