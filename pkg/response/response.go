package response

import (
	"github.com/labstack/echo/v4"
)

type Success struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Error struct {
	Message string `json:"message"`
}

func SuccessResponse(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, Success{
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, status int, message string) error {
	return c.JSON(status, Error{
		Message: message,
	})
}
