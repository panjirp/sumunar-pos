package dto

import "time"

type UserResponse struct {
	ID        string     `json:"id"`
	Fullname  string     `json:"fullname"`
	Email     string     `json:"email"`
	GoogleID  *string    `json:"google_id,omitempty"`
	Picture   *string    `json:"picture,omitempty"`
	Provider  string     `json:"provider"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
