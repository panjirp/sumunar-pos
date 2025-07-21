package dto

type StoreResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Phone     *string `json:"phone,omitempty"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	CreatedBy string  `json:"created_by"`
	UpdatedAt string  `json:"updated_at"`
	UpdatedBy string  `json:"updated_by"`
}
