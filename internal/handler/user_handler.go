package handler

import (
	"errors"
	"net/http"
	"strconv"

	"rentalin/internal/dto"
	"rentalin/internal/errs"
	"rentalin/internal/service"
	"rentalin/pkg/response"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(int)

	ctx := c.Request().Context()

	profile, err := h.userService.GetProfile(
		ctx,
		userID,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, "User not found")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Profile fetched", profile)
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.userService.GetAllUsers(
		ctx,
	)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
	}

	return response.SuccessResponse(c, http.StatusOK, "Users fetched", users)
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	ctx := c.Request().Context()

	err := h.userService.UpdateProfile(
		ctx,
		userID,
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput),
			errors.Is(err, errs.ErrInvalidUsernameFormat),
			errors.Is(err, errs.ErrInvalidPhoneFormat),
			errors.Is(err, errs.ErrInvalidAddress),
			errors.Is(err, errs.ErrUsernameContainsSpace):

			return response.ErrorResponse(c, http.StatusBadRequest, err.Error())

		case errors.Is(err, errs.ErrUserNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, "User not found")

		case errors.Is(err, errs.ErrUserUsernameExists):
			return response.ErrorResponse(c, http.StatusConflict, "Username already exists")

		case errors.Is(err, errs.ErrUserPhoneExists):
			return response.ErrorResponse(c, http.StatusConflict, "Phone already exists")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Profile updated", nil)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request parameter")
	}

	ctx := c.Request().Context()

	err = h.userService.DeleteUser(
		ctx,
		userID,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, "User not found")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "User deleted", nil)
}
