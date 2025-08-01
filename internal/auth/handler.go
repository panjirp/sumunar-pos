package auth

import (
	"fmt"
	"net/http"

	"sumunar-pos-core/internal/auth/dto"
	"sumunar-pos-core/middleware"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authService Service
}

func NewHandler(as Service) *Handler {
	return &Handler{authService: as}
}

// Login godoc
// @Summary Login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login data"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, refreshToken, u, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	resp := dto.AuthResponse{
		ID:           u.ID,
		Fullname:     u.Fullname,
		Email:        u.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": resp,
	})
}

// Register godoc
// @Summary Register new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register data"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Register user
	token, refreshToken, u, err := h.authService.RegisterManual(c.Request().Context(), req)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	resp := dto.AuthResponse{
		ID:           u.ID,
		Fullname:     u.Fullname,
		Email:        u.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"data": resp,
	})
}

// GoogleLogin godoc
// @Summary Login using Google ID Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.GoogleLoginRequest true "Google Login"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/google [post]
func (h *Handler) GoogleLogin(c echo.Context) error {
	var req dto.GoogleLoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Lanjutkan ke login/daftar user
	user, err := h.authService.LoginWithGoogle(c.Request().Context(), req.IDToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Google login failed")
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Token generation failed")
	}

	return c.JSON(http.StatusOK, dto.AuthResponse{
		ID:       user.ID,
		Fullname: user.Fullname,
		Email:    user.Email,
		Token:    token,
	})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh Token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c echo.Context) error {
	var req dto.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, newRefreshToken, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

// Logout godoc
// @Summary Logout and revoke refresh token
// @Tags Auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param refresh_token formData string true "Refresh Token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *Handler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(string)
	refreshToken := c.FormValue("refresh_token")

	err := h.authService.Logout(c.Request().Context(), refreshToken, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully logged out",
	})
}
