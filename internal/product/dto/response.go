package dto

type ProductResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	StoreID string `json:"store_id"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
