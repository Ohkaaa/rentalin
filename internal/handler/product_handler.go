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

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req dto.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	ctx := c.Request().Context()

	product, err := h.productService.CreateProduct(
		ctx,
		req.Name,
		req.DailyPrice,
		req.Stock,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput),
			errors.Is(err, errs.ErrInvalidPrice),
			errors.Is(err, errs.ErrInvalidStock):

			return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Product created successfully", product)
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid product id")
	}

	ctx := c.Request().Context()

	product, err := h.productService.GetProductByID(
		ctx,
		productID,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrProductNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, "Product not found")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Product fetched", product)
}

func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.productService.GetAllProducts(
		ctx,
	)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
	}

	return response.SuccessResponse(c, http.StatusOK, "Products fetched", products)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid product id")
	}

	var req dto.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	ctx := c.Request().Context()

	err = h.productService.UpdateProduct(
		ctx,
		productID,
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput),
			errors.Is(err, errs.ErrInvalidPrice),
			errors.Is(err, errs.ErrInvalidStock):

			return response.ErrorResponse(c, http.StatusBadRequest, err.Error())

		case errors.Is(err, errs.ErrProductNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, "Product not found")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Product updated", nil)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid product id")
	}

	ctx := c.Request().Context()

	err = h.productService.DeleteProduct(
		ctx,
		productID,
	)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrProductNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, "Product not found")

		default:
			return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	return response.SuccessResponse(c, http.StatusOK, "Product deleted", nil)
}
