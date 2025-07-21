package dto

type UserRequest struct {
	Username string  `json:"username" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password *string `json:"password,omitempty"` // hanya untuk login/register manual
	GoogleID *string `json:"google_id,omitempty"`
	Picture  *string `json:"picture,omitempty"`
	Provider string  `json:"provider" validate:"required,oneof=local google"`
	Role     string  `json:"role" validate:"required,oneof=owner worker"`
}
