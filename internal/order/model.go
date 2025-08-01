package order

import "time"

type Order struct {
	ID            string    `json:"id"`
	StoreID       string    `json:"store_id"`
	InvoiceNumber string    `json:"invoice_number"`
	CustomerID    string    `json:"customer_id"`
	Status        string    `json:"status"`   // pending, processed, done, taken
	Discount      float64   `json:"discount"` // e.g. 10.00 for 10%
	TotalPrice    float64   `json:"total_price"`
	PaidAmount    float64   `json:"paid_amount"`
	Change        float64   `json:"change"`
	PickupDate    time.Time `json:"pickup_date"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"` // user ID dari kasir
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
}

type OrderItem struct {
	ID               string  `json:"id"`
	OrderID          string  `json:"order_id"`
	ProductServiceID string  `json:"product_service_id"`
	Quantity         float64 `json:"quantity"` // 3 pcs, 2.5 kg, etc
	TotalPrice       float64 `json:"total_price"`
	Notes            string  `json:"notes"` // opsional, misal "khusus pakaian putih"
}
