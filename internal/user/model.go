package user

import (
	"sumunar-pos-core/internal/base"
	"time"
)

type User struct {
	ID        string     `json:"id"`
	Fullname  string     `json:"fullname"`
	Email     string     `json:"email"`
	GoogleID  *string    `json:"google_id,omitempty"` // nullable jika manual
	Password  *string    `json:"password,omitempty"`  // nullable jika pakai Google
	Picture   *string    `json:"picture,omitempty"`   // avatar dari Google
	Provider  string     `json:"provider"`            // "google" atau "local"
	LastLogin *time.Time `json:"last_login,omitempty"`
	Role      string     `json:"role"` // "owner", "worker"
	base.BaseModel
}
