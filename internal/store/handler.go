package store

import (
	"net/http"
	"strconv"

	"sumunar-pos-core/internal/store/dto"

	"github.com/labstack/echo/v4"
)

type StoreHandler struct {
	service StoreService
}

func NewHandler(service StoreService) *StoreHandler {
	return &StoreHandler{service}
}

func (h *StoreHandler) Create(c echo.Context) error {
	var req dto.StoreRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	store, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, ToStoreResponse(store))
}

func (h *StoreHandler) FindByID(c echo.Context) error {
	id := c.Param("id")

	store, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, ToStoreResponse(store))
}

func (h *StoreHandler) FindAll(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	stores, total, err := h.service.FindAll(c.Request().Context(), limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responses := ToStoreListResponse(stores)

	return c.JSON(http.StatusOK, echo.Map{
		"data":   responses,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *StoreHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req dto.StoreRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	store, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ToStoreResponse(store))
}

func (h *StoreHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
