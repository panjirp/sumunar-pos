package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"sumunar-pos-core/internal/auth/dto"
	"sumunar-pos-core/internal/user"
	"sumunar-pos-core/middleware"
	"sumunar-pos-core/pkg/hash"
	"sumunar-pos-core/pkg/utils"

	"github.com/labstack/echo/v4"

	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type Service interface {
	RegisterManual(ctx context.Context, req dto.RegisterRequest) (string, string, *user.User, error)
	Login(ctx context.Context, email, password string) (string, string, *user.User, error)
	LoginWithGoogle(ctx context.Context, IDToken string) (*user.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken, userID string) error
}

type service struct {
	userRepo         user.UserRepository
	refreshTokenRepo RefreshTokenRepository
}

func NewService(userRepo user.UserRepository, refreshTokenRepo RefreshTokenRepository) Service {
	return &service{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}

func (s *service) RegisterManual(ctx context.Context, req dto.RegisterRequest) (string, string, *user.User, error) {

	// Check if user exists (by email)
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusBadRequest, "Email already registered")
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	u := ToUserModel(&req, hashedPassword)

	token, err := middleware.GenerateJWT(u.ID, u.Role)
	if err != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = s.userRepo.Create(ctx, u)
	if err != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	refresh := ToRefreshTokenModel(u.ID)
	err = s.refreshTokenRepo.Create(ctx, refresh)
	if err != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return token, refresh.Token, u, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, string, *user.User, error) {
	// Verifikasi email & password
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || !hash.CheckPasswordHash(password, *u.Password) {
		return "", "", nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	// update last_login
	_ = s.userRepo.UpdateLastLogin(ctx, u.ID)

	token, err := middleware.GenerateJWT(u.ID, u.Role)
	if err != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	refresh := ToRefreshTokenModel(u.ID)
	err = s.refreshTokenRepo.RevokeAllByUser(ctx, u.ID)
	if err != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = s.refreshTokenRepo.Create(ctx, refresh)
	if err != nil {
		return "", "", nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return token, refresh.Token, u, nil
}

func (s *service) LoginWithGoogle(ctx context.Context, IDToken string) (*user.User, error) {

	// ✅ Verifikasi token dari Google
	payload, err := idtoken.Validate(context.Background(), IDToken, "")
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid Google token")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Email not found in Google token")
	}

	name, _ := payload.Claims["name"].(string)
	sub, _ := payload.Claims["sub"].(string)

	// Cek apakah user sudah ada
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && err != user.ErrUserNotFound {
		return nil, err
	}

	if existingUser != nil {
		// Update last login
		existingUser.LastLogin = utils.PtrTime(time.Now())
		if err := s.userRepo.UpdateLastLogin(ctx, existingUser.ID); err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	// Jika belum ada → daftar
	newUser := &user.User{
		ID:        uuid.New().String(),
		Fullname:  name,
		Email:     email,
		GoogleID:  &sub,
		Role:      "user",
		LastLogin: utils.PtrTime(time.Now()),
	}

	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	t, err := s.refreshTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired refresh token")
	}

	user, err := s.userRepo.FindByID(ctx, t.UserID)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	accessToken, err := middleware.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate JWT")
	}

	// revoke used refresh token and generate new one
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	newRefresh := ToRefreshTokenModel(userID)
	_ = s.refreshTokenRepo.Revoke(ctx, t.ID)
	err = s.refreshTokenRepo.Create(ctx, newRefresh)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return accessToken, newRefresh.Token, nil
}

func (s *service) Logout(ctx context.Context, refreshToken, userID string) error {

	t, err := s.refreshTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil || t.UserID != userID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	if err := s.refreshTokenRepo.Revoke(ctx, t.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke token")
	}

	return nil
}
