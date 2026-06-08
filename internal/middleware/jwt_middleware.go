package middleware

import (
	"net/http"
	"rentalin/pkg/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "missing token"})
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			claims, err := auth.ParseJWT(tokenString, secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token"})
			}

			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}
