package handler

import (
	"errors"
	"net/http"
	"rentalin/internal/dto"
	"rentalin/internal/errs"
	"rentalin/internal/service"
	"rentalin/pkg/response"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	ctx := c.Request().Context()

	user, err := h.authService.Register(
		ctx,
		req.Username,
		req.Email,
		req.Phone,
		req.Address,
		req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput),
			errors.Is(err, errs.ErrInvalidAddress),
			errors.Is(err, errs.ErrInvalidUsernameFormat),
			errors.Is(err, errs.ErrInvalidEmailFormat),
			errors.Is(err, errs.ErrInvalidPhoneFormat),
			errors.Is(err, errs.ErrUsernameContainsSpace),
			errors.Is(err, errs.ErrWeakPassword):

			return response.ErrorResponse(c, http.StatusBadRequest, err.Error())

		case errors.Is(err, errs.ErrUserUsernameExists):
			return response.ErrorResponse(c, http.StatusConflict, "Username already registered")

		case errors.Is(err, errs.ErrUserEmailExists):
			return response.ErrorResponse(c, http.StatusConflict, "Email already registered")

		case errors.Is(err, errs.ErrUserPhoneExists):
			return response.ErrorResponse(c, http.StatusConflict, "Phone number already registered")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	ctx := c.Request().Context()

	token, err := h.authService.Login(
		ctx,
		req.Email,
		req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return response.ErrorResponse(c, http.StatusBadRequest, err.Error())

		case errors.Is(err, errs.ErrInvalidCredentials):
			return response.ErrorResponse(c, http.StatusUnauthorized, "Email or password is incorrect")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Login successful", map[string]string{
		"token": token,
	})
}
