package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"rentalin/internal/dto"
	"rentalin/internal/errs"
	"rentalin/internal/service"
	"rentalin/pkg/response"

	"github.com/labstack/echo/v4"
)

type RentalHandler struct {
	rentalService service.RentalService
}

func NewRentalHandler(rentalService service.RentalService) *RentalHandler {
	return &RentalHandler{
		rentalService: rentalService,
	}
}

func (h *RentalHandler) CreateRental(c echo.Context) error {
	var req dto.CreateRentalRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	ctx := c.Request().Context()

	userID := c.Get("user_id").(int)
	role := c.Get("role").(string)

	customerID := userID

	if role == "admin" {
		if req.CustomerID != nil {
			customerID = *req.CustomerID
		} else {
			return response.ErrorResponse(c, http.StatusBadRequest, "customer_id is required for admin")
		}
	}

	rental, err := h.rentalService.CreateRental(
		ctx,
		customerID,
		req.ProductID,
		userID,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput),
			errors.Is(err, errs.ErrInvalidDate),
			errors.Is(err, errs.ErrInvalidRentalPeriod):

			return response.ErrorResponse(c, http.StatusBadRequest, err.Error())

		case errors.Is(err, errs.ErrProductOutOfStock):
			return response.ErrorResponse(c, http.StatusConflict, "Product is out of stock")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Rental created successfully", rental)
}

func (h *RentalHandler) GetRentalByID(c echo.Context) error {
	userID := c.Get("user_id").(int)
	role := c.Get("role").(string)

	rentalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid rental ID")
	}

	ctx := c.Request().Context()

	rental, err := h.rentalService.GetRentalByID(
		ctx,
		rentalID,
		userID,
		role,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRentalNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, "Rental not found")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Rental fetched", rental)
}

func (h *RentalHandler) GetRentalsByCustomerID(c echo.Context) error {
	userID := c.Get("user_id").(int)

	ctx := c.Request().Context()

	rentals, err := h.rentalService.GetRentalsByCustomerID(
		ctx,
		userID,
	)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
	}

	return response.SuccessResponse(c, http.StatusOK, "Rentals fetched", rentals)

}

func (h *RentalHandler) GetAllRentals(c echo.Context) error {
	ctx := c.Request().Context()

	rentals, err := h.rentalService.GetAllRentals(
		ctx,
	)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
	}

	return response.SuccessResponse(c, http.StatusOK, "Rentals fetched", rentals)
}

func (h *RentalHandler) CancelRental(c echo.Context) error {
	rentalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request parameter")
	}

	userID := c.Get("user_id").(int)
	role := c.Get("role").(string)

	ctx := c.Request().Context()

	err = h.rentalService.CancelRental(
		ctx,
		rentalID,
		userID,
		role,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRentalNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, "Rental not found")

		case errors.Is(err, errs.ErrForbidden):
			return response.ErrorResponse(c, http.StatusForbidden, "Forbidden")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Rental cancelled", nil)
}

func (h *RentalHandler) CompleteRental(c echo.Context) error {
	rentalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request parameter")
	}

	ctx := c.Request().Context()

	err = h.rentalService.CompleteRental(
		ctx,
		rentalID,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return response.ErrorResponse(c, http.StatusBadRequest, "Cannot complete rental")

		default:
			log.Println(err)

			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Rental completed", nil)
}
