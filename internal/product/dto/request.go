package dto

type ProductRequest struct {
	Name    string `json:"name" validate:"required"`
	StoreID string `json:"store_id" validate:"required"`
}
