package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			claims, ok := c.Get("user").(Claims)
			if !ok {
				return c.JSON(
					http.StatusUnauthorized,
					map[string]string{
						"message": "authentication required",
					},
				)
			}

			for _, userRole := range claims.RealmAccess.Roles {
				if userRole == role {
					return next(c)
				}
			}

			return c.JSON(
				http.StatusForbidden,
				map[string]string{
					"message": "insufficient permissions",
				},
			)
		}
	}
}
