package dto

type AuthResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// type RefreshTokenResponse struct {
// 	ID        uuid.UUID
// 	Token     string
// 	UserID    string
// 	ExpiresAt time.Time
// 	Revoked   bool
// 	CreatedAt time.Time
// }
