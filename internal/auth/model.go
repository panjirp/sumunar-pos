package auth

import (
	"sumunar-pos-core/internal/base"
	"time"
)

type RefreshToken struct {
	ID        string    `json:"id"`         // UUID
	UserID    string    `json:"user_id"`    // FK ke user table
	Token     string    `json:"token"`      // Refresh token string (biasanya UUID atau panjang acak)
	ExpiresAt time.Time `json:"expires_at"` // Expiration time
	Revoked   bool      `json:"revoked"`    // Untuk revoke manual (logout, banned)
	base.BaseModel
}
