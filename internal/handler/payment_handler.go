package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"rentalin/config"
	"rentalin/internal/dto"
	"rentalin/internal/errs"
	"rentalin/internal/service"
	"rentalin/pkg/response"

	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	cfg            *config.Config
	paymentService service.PaymentService
}

func NewPaymentHandler(cfg *config.Config, paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		cfg:            cfg,
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	var req dto.CreatePaymentRequest

	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	userID, ok := c.Get("user_id").(int)
	if !ok {
		return response.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	ctx := c.Request().Context()

	payment, err := h.paymentService.CreatePayment(
		ctx,
		req.RentalID,
		userID,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRentalNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, "Rental not found")

		case errors.Is(err, errs.ErrForbidden):
			return response.ErrorResponse(c, http.StatusForbidden, "Forbidden")

		case errors.Is(err, errs.ErrInvalidInput):
			return response.ErrorResponse(c, http.StatusBadRequest, "Invalid input")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusCreated, "Payment created", payment)
}

func (h *PaymentHandler) GetPaymentsByCustomerID(c echo.Context) error {
	userID := c.Get("user_id").(int)

	ctx := c.Request().Context()

	payments, err := h.paymentService.GetPaymentsByCustomerID(
		ctx,
		userID,
	)
	if err != nil {
		log.Println(err)
		return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
	}

	return response.SuccessResponse(c, http.StatusOK, "Payments fetched", payments)
}

func (h *PaymentHandler) GetAllPayments(c echo.Context) error {
	ctx := c.Request().Context()

	payments, err := h.paymentService.GetAllPayments(
		ctx,
	)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
	}

	return response.SuccessResponse(c, http.StatusOK, "Payments fetched", payments)
}

func (h *PaymentHandler) XenditWebhook(c echo.Context) error {
	callbackToken := c.Request().Header.Get("X-CALLBACK-TOKEN")

	if callbackToken != h.cfg.XenditCallbackToken {
		return response.ErrorResponse(c, http.StatusUnauthorized, "Invalid callback token")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Failed to read request body")
	}

	var payload dto.XenditWebhookPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
	}

	ctx := c.Request().Context()

	switch payload.Status {
	case "PAID":
		err = h.paymentService.HandleInvoicePaid(
			ctx,
			payload.ExternalID,
			payload.PaidAmount,
			payload.PaymentMethod,
			payload.PaymentChannel,
			body,
		)

	case "=EXPIRED":
		err = h.paymentService.HandleInvoiceExpired(
			ctx,
			payload.ExternalID,
			body,
		)

	case "FAILED":
		err = h.paymentService.HandleInvoiceFailed(
			ctx,
			payload.ExternalID,
			body,
		)

	default:
		return response.SuccessResponse(c, http.StatusOK, "Ignored", nil)
	}

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Failed to process payment")
	}

	return response.SuccessResponse(c, http.StatusOK, "Success", nil)
}
