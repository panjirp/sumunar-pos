package servicetype

import (
	"net/http"
	"strconv"

	"sumunar-pos-core/internal/servicetype/dto"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service ServiceTypeService
}

func NewHandler(service ServiceTypeService) *Handler {
	return &Handler{service}
}

// Create godoc
// @Summary Create a new service type
// @Tags servicetypes
// @Accept json
// @Produce json
// @Param request body dto.StoreRequest true "Store request"
// @Success 201 {object} dto.StoreResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /servicetypes [post]
func (h *Handler) Create(c echo.Context) error {
	var req dto.ServiceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	service, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, ToServiceTypeResponse(service))
}

func (h *Handler) FindByID(c echo.Context) error {
	id := c.Param("id")

	service, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, ToServiceTypeResponse(service))
}

func (h *Handler) FindAll(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	services, total, err := h.service.FindAll(c.Request().Context(), limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responses := ToServiceTypeListResponse(services)

	return c.JSON(http.StatusOK, echo.Map{
		"data":   responses,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var req dto.ServiceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	service, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ToServiceTypeResponse(service))
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
