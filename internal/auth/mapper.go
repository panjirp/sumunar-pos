package auth

import (
	"time"

	"sumunar-pos-core/internal/auth/dto"
	"sumunar-pos-core/internal/base"
	"sumunar-pos-core/internal/user"

	"sumunar-pos-core/pkg/utils"

	"github.com/google/uuid"
)

func ToUserModel(req *dto.RegisterRequest, hashedPassword string) *user.User {
	now := time.Now()
	return &user.User{
		ID:        uuid.New().String(),
		Fullname:  req.Fullname,
		Email:     req.Email,
		Password:  &hashedPassword,
		Provider:  "local",
		Role:      "owner",
		LastLogin: &now,
		BaseModel: base.BaseModel{
			IsActive:  true,
			CreatedAt: now,
			CreatedBy: "SYSTEM",
			UpdatedAt: now,
			UpdatedBy: "SYSTEM",
		},
	}
}

func ToRefreshTokenModel(userID string) *RefreshToken {
	now := time.Now()
	return &RefreshToken{
		ID:        uuid.New().String(),
		Token:     utils.GenerateRandomToken(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		BaseModel: base.BaseModel{
			IsActive:  true,
			CreatedAt: now,
			CreatedBy: userID,
			UpdatedAt: now,
			UpdatedBy: userID,
		},
	}
}
