package auth

import (
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

// login
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
		Username:     u.Username,
		Email:        u.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": resp,
	})
}

// register
func (h *Handler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Register user
	u, err := h.authService.RegisterManual(c.Request().Context(), req.Username, req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Registration failed")
	}

	resp := dto.AuthResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Token:    "",
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"data": resp,
	})
}

// login with google
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
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	})
}

// refresh token
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
