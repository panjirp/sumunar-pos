package dto

type ProductServiceRequest struct {
	ProductID     string  `json:"product_id" validate:"required"`
	ServiceTypeID string  `json:"service_type_id" validate:"required"`
	Unit          string  `json:"unit" validate:"required"`
	Price         float64 `json:"price" validate:"required"`
}
