package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		role, ok := c.Get("role").(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		}

		if role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "forbidden"})
		}

		return next(c)
	}
}
