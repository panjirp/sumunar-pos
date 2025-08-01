package dto

type OrderResponse struct {
	ID              string              `json:"id"`
	InvoiceNumber   string              `json:"invoice_number"`
	StoreID         string              `json:"store_id"`
	CustomerID      string              `json:"customer_id"`
	CustomerName    string              `json:"customer_name"`
	CustomerPhone   string              `json:"customer_phone"`
	CustomerAddress string              `json:"customer_address"`
	Status          string              `json:"status"`
	Discount        float64             `json:"discount"`
	TotalPrice      float64             `json:"total_price"`
	PaidAmount      float64             `json:"paid_amount"`
	Change          float64             `json:"change"`
	PickupDate      string              `json:"pickup_date"` // ISO8601
	CreatedAt       string              `json:"created_at"`
	CreatedBy       string              `json:"created_by"`
	OrderItems      []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ID               string  `json:"id"`
	ProductServiceID string  `json:"product_service_id"`
	Quantity         float64 `json:"quantity"`
	TotalPrice       float64 `json:"total_price"`
	Notes            string  `json:"notes"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
