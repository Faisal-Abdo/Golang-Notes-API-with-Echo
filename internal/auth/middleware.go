package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	provider *oidc.Provider        // OIDC provider for Keycloak
	verifier *oidc.IDTokenVerifier // Verifier for validating JWT tokens
}

func NewAuthMiddleware(keycloakURL, realm, clientID string) (*AuthMiddleware, error) {
	ctx := context.Background()

	issuerURL := fmt.Sprintf("%s/realms/%s", keycloakURL, realm)
	provider, err := oidc.NewProvider(ctx, issuerURL)

	if err != nil {
		return nil, fmt.Errorf("failed to create oidc provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
		// Keycloak sets aud to "account" by default, not the client ID, so the
		// default audience check would always fail here.
		SkipClientIDCheck: true,
	})

	return &AuthMiddleware{provider: provider, verifier: verifier}, nil
}

func extractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header must start with Bearer")
	}
	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", errors.New("bearer token is empty")
	}
	return token, nil
}

func (a *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rawToken, err := extractBearerToken(c.Request().Header.Get("Authorization"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

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
