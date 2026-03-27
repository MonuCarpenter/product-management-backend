package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func RoleMiddleware(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("userRole")
			if userRole == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			roleStr := userRole.(string)
			for _, r := range roles {
				if strings.EqualFold(r, roleStr) {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
	}
}
