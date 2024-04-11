package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var ( // env keys
	envKeyDatabaseURL              = "DB_URL"
	envKeySigningKeyRegistration   = "SIGNING_KEY_REGISTRATION"
	envKeySigningKeyAuthentication = "SIGNING_KEY_AUTHENTICATION"
	envKeyHTTPPort                 = "HTTP_PORT"
	envKeyGRPCPort                 = "GRPC_PORT"
	envKeyGrpcTLSCertFile          = "GRPC_TLS_CERT_FILE"
	envKeyGrpcTLSKeyFile           = "GRPC_TLS_KEY_FILE"
	envKeyHTTPTLSKeyFile           = "HTTP_TLS_KEY_FILE"
	envKeyHTTPTLSCertFile          = "HTTP_TLS_CERT_FILE"
	envKeyRelayUrl                 = "RELAY_URL"
	envKeyAuthUrl                  = "AUTH_URL"
)

// GetConfigString returns the value of the key as a string.
func GetConfigString(key string) (string, error) {
	if !viper.IsSet(key) {
		return "", fmt.Errorf("%s is not set", key)
	}

	value := viper.GetString(key)
	if value == "" {
		return "", fmt.Errorf("%s is empty", key)
	}

	return value, nil
}

// GetConfigInt returns the value of the key as an integer.
func GetConfigInt(key string) (int, error) {
	if !viper.IsSet(key) {
		return 0, fmt.Errorf("%s is not set", key)
	}

	return viper.GetInt(key), nil
}

func DatabaseURL() (string, error) {
	return GetConfigString(envKeyDatabaseURL)
}

func SigningKeyRegistration() (string, error) {
	return GetConfigString(envKeySigningKeyRegistration)
}

func SigningKeyAuthentication() (string, error) {
	return GetConfigString(envKeySigningKeyAuthentication)
}

func HTTPPort() (int, error) {
	return GetConfigInt(envKeyHTTPPort)
}

func GRPCPort() (int, error) {
	return GetConfigInt(envKeyGRPCPort)
}

func GrpcTLSCertFile() (string, error) {
	return GetConfigString(envKeyGrpcTLSCertFile)
}

func GrpcTLSKeyFile() (string, error) {
	return GetConfigString(envKeyGrpcTLSKeyFile)
}

func HTTPTLSKeyFile() (string, error) {
	return GetConfigString(envKeyHTTPTLSKeyFile)
}

func HTTPTLSCertFile() (string, error) {
	return GetConfigString(envKeyHTTPTLSCertFile)
}

func RelayUrl() (string, error) {
	return GetConfigString(envKeyRelayUrl)
}

func AuthUrl() (string, error) {
	return GetConfigString(envKeyAuthUrl)
}
