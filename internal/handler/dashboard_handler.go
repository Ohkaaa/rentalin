package handler

import (
	"net/http"
	"rentalin/internal/service"
	"rentalin/pkg/response"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetDashboard(c echo.Context) error {
	ctx := c.Request().Context()

	info, err := h.dashboardService.GetDashboard(
		ctx,
	)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")

	}

	return response.SuccessResponse(c, http.StatusOK, "Dashboard fetched", info)
}
