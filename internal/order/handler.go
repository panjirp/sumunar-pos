package order

import (
	"net/http"
	"strconv"

	"sumunar-pos-core/internal/order/dto"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *OrderService
}

func NewHandler(svc *OrderService) *Handler {
	return &Handler{service: svc}
}

// Create order
func (h *Handler) Create(c echo.Context) error {
	var req dto.OrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}

	userID := c.Get("user_id").(string)

	resp, err := h.service.CreateOrder(c.Request().Context(), req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

// FindAll orders
func (h *Handler) FindAll(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit <= 0 {
		limit = 20
	}

	orders, total, err := h.service.FindAll(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data":  orders,
		"total": total,
	})
}

// GetByID
// func (h *Handler) FindByID(c echo.Context) error {
// 	id := c.Param("id")

// 	order, err := h.service.FindByID(c.Request().Context(), id)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Order not found"})
// 	}

// 	return c.JSON(http.StatusOK, order)
// }

// Update
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var req dto.OrderRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
	}

	userID := c.Get("user_id").(string)

	resp, err := h.service.Update(c.Request().Context(), id, &req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// Delete
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
