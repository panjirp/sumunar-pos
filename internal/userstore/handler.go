package userstore

import (
	"net/http"

	"sumunar-pos-core/middleware"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type AssignRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	StoreID string `json:"store_id" validate:"required"`
}

func (h *Handler) Assign(c echo.Context) error {
	var req AssignRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := h.service.AssignUserToStore(c.Request().Context(), req.UserID, req.StoreID, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User assigned to store"})
}
