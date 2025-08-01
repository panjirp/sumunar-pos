package dto

type OrderRequest struct {
	StoreID         string             `json:"store_id" validate:"required"`
	CustomerID      string             `json:"customer_id"`
	CustomerName    string             `json:"customer_name"`
	CustomerPhone   string             `json:"customer_phone"`
	CustomerAddress string             `json:"customer_address"`
	Status          string             `json:"status" validate:"required,oneof='PENDING' 'ONGOING' 'COMPLETED' 'CANCELLED'"`
	PickupDate      string             `json:"pickup_date"` // ISO8601
	Items           []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
	Discount        float64            `json:"discount"`    // persen: 0 - 100
	PaidAmount      float64            `json:"paid_amount"` // nilai yang dibayar
}

type OrderItemRequest struct {
	ProductServiceID string  `json:"product_service_id" validate:"required"`
	Quantity         float64 `json:"quantity" validate:"required,gt=0"`
	TotalPrice       float64 `json:"total_price"`
	Notes            string  `json:"notes"` // optional
}
