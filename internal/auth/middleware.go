package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	provider *oidc.Provider        // OIDC provider for Keycloak
	verifier *oidc.IDTokenVerifier // Verifier for validating JWT tokens
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

	issuerURL := fmt.Sprintf("%s/realms/%s", keycloakURL, realm)
	provider, err := oidc.NewProvider(ctx, issuerURL)

	if err != nil {
		return nil, fmt.Errorf("failed to create oidc provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: os.Getenv("KEYCLOAK_CLIENT_ID"),
	})

	return &AuthMiddleware{provider: provider, verifier: verifier}, nil
}

func (a *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		return next(c)
	}
}
