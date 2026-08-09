package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL      string
	KeycloakURL      string
	KeycloakRealm    string
	KeycloakClientID string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		KeycloakURL:      os.Getenv("KEYCLOAK_URL"),
		KeycloakRealm:    os.Getenv("KEYCLOAK_REALM"),
		KeycloakClientID: os.Getenv("KEYCLOAK_CLIENT_ID"),
	}

	var missing []string
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if cfg.KeycloakURL == "" {
		missing = append(missing, "KEYCLOAK_URL")
	}
	if cfg.KeycloakRealm == "" {
		missing = append(missing, "KEYCLOAK_REALM")
	}
	if cfg.KeycloakClientID == "" {
		missing = append(missing, "KEYCLOAK_CLIENT_ID")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}
