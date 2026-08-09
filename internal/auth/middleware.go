package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

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
		// Keycloak access tokens set aud to "account" by default, not the
		// client ID, so the default audience check would always fail here.
		// The client is validated via the azp claim instead, below.
		SkipClientIDCheck: true,
	})

	return &AuthMiddleware{provider: provider, verifier: verifier}, nil
}

func (a *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is missing")
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header must start with Bearer")
		}
		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		idToken, err := a.verifier.Verify(c.Request().Context(), rawToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}
		var claims Claims

		if err := idToken.Claims(&claims); err != nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"invalid token claims",
			)
		}
		c.Set("user", claims)
		return next(c)
	}
}
