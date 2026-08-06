package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
)

type AuthMiddleware struct {
	provider *oidc.Provider
}

func NewAuthMiddleware() (*AuthMiddleware, error) {
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	if keycloakURL == "" {
		return nil, fmt.Errorf("KEYCLOAK_URL environment variable is not set")
	}
	realm := os.Getenv("KEYCLOAK_REALM")
	if realm == "" {
		return nil, fmt.Errorf("KEYCLOAK_REALM environment variable is not set")
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, keycloakURL+"/realms/"+realm)
	if err != nil {
		return nil, fmt.Errorf("failed to create oidc provider: %w", err)
	}

	return &AuthMiddleware{provider: provider}, nil
}
