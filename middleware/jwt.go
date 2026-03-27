package middleware

import (
	"net/http"
	"strings"
	"product-management-backend/auth"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid token"})
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ParseJWT(tokenStr)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}
		c.Set("userId", claims.UserID)
		c.Set("userRole", claims.Role)
		return next(c)
	}
}
