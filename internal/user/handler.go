package user

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"sumunar-pos-core/internal/user/dto"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) ListUsers(c echo.Context) error {
	ctx := c.Request().Context()
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	page := (offset / limit) + 1

	users, total, err := h.service.ListUsers(ctx, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// pages := int(math.Ceil(float64(total) / float64(limit)))
	res := NewPaginatedResponse(users, total, page, limit)

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.UserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	user, err := h.service.CreateUser(ctx, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}
