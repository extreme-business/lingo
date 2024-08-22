package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const ( // env keys.
	databaseURL               = "DB_URL"
	signingKeyAccessToken     = "SIGNING_KEY_ACCESS_TOKEN"
	signingKeyRefreshToken    = "SIGNING_KEY_REFRESH_TOKEN"
	httpPort                  = "HTTP_PORT"
	grpcPort                  = "GRPC_PORT"
	accountServiceTLSCertFile = "ACCOUNT_TLS_CERT_FILE"
	accountServiceTLSKeyFile  = "ACCOUNT_TLS_KEY_FILE"
	httpTLSKeyFile            = "HTTP_TLS_KEY_FILE"
	httpTLSCertFile           = "HTTP_TLS_CERT_FILE"
	accountServiceURL         = "ACCOUNT_SERVICE_URL"
	sysUserID                 = "SYSTEM_USER_ID"
	sysUserEmail              = "SYSTEM_USER_EMAIL"
	sysUserPassword           = "SYSTEM_USER_PASSWORD"
	sysOrganizationID         = "SYSTEM_ORGANIZATION_ID"
	sysOrgLegalName           = "SYSTEM_ORGANIZATION_LEGAL_NAME"
	sysOrgSlug                = "SYSTEM_ORGANIZATION_SLUG"
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

// str returns the value of the key as a string.
func (c *Config) str(key string) (string, error) {
	if !c.viper.IsSet(key) {
		return "", fmt.Errorf("%s is not set", key)
	}

	return c.viper.GetString(key), nil
}

// int returns the value of the key as an integer.
func (c *Config) int(key string) (int, error) {
	if !c.viper.IsSet(key) {
		return 0, fmt.Errorf("%s is not set", key)
	}

	return c.viper.GetInt(key), nil
}

func (c *Config) DatabaseURL() (string, error)                 { return c.str(databaseURL) }
func (c *Config) SigningKeyAccessToken() (string, error)       { return c.str(signingKeyAccessToken) }
func (c *Config) SigningKeyRefreshToken() (string, error)      { return c.str(signingKeyRefreshToken) }
func (c *Config) HTTPPort() (int, error)                       { return c.int(httpPort) }
func (c *Config) GRPCPort() (int, error)                       { return c.int(grpcPort) }
func (c *Config) AccountServiceTLSCertFile() (string, error)   { return c.str(accountServiceTLSCertFile) }
func (c *Config) AccountServiceTLSKeyFile() (string, error)    { return c.str(accountServiceTLSKeyFile) }
func (c *Config) HTTPTLSKeyFile() (string, error)              { return c.str(httpTLSKeyFile) }
func (c *Config) HTTPTLSCertFile() (string, error)             { return c.str(httpTLSCertFile) }
func (c *Config) AccountServiceURL() (string, error)           { return c.str(accountServiceURL) }
func (c *Config) SystemUserID() (string, error)                { return c.str(sysUserID) }
func (c *Config) SystemUserEmail() (string, error)             { return c.str(sysUserEmail) }
func (c *Config) SystemUserPassword() (string, error)          { return c.str(sysUserPassword) }
func (c *Config) SystemOrganizationID() (string, error)        { return c.str(sysOrganizationID) }
func (c *Config) SystemOrganizationLegalName() (string, error) { return c.str(sysOrgLegalName) }
func (c *Config) SystemOrganizationSlug() (string, error)      { return c.str(sysOrgSlug) }
