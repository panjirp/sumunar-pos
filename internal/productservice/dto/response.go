package dto

type ProductServiceResponse struct {
	ID            string  `json:"id"`
	ProductID     string  `json:"product_id"`
	ServiceTypeID string  `json:"service_type_id"`
	Unit          string  `json:"unit"`
	Price         float64 `json:"price"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
