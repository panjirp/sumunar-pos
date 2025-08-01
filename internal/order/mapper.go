package order

import (
	"time"

	"sumunar-pos-core/internal/customer"
	"sumunar-pos-core/internal/order/dto"

	"github.com/google/uuid"
)

// Request → Model
func ToOrderModel(req dto.OrderRequest, createdBy string) (*Order, []*OrderItem, error) {
	orderID := uuid.New().String()
	now := time.Now()

	order := &Order{
		ID:            orderID,
		StoreID:       req.StoreID,
		InvoiceNumber: "", // akan digenerate di service
		CustomerID:    req.CustomerID,
		Status:        "pending",
		Discount:      req.Discount,
		TotalPrice:    0, // akan dihitung di service
		PaidAmount:    req.PaidAmount,
		Change:        0, // akan dihitung di service
		PickupDate:    mustParseTime(req.PickupDate),
		CreatedAt:     now,
		CreatedBy:     createdBy,
		UpdatedAt:     now,
		UpdatedBy:     createdBy,
	}

	var items []*OrderItem
	for _, item := range req.Items {
		items = append(items, &OrderItem{
			ID:               uuid.New().String(),
			OrderID:          orderID,
			ProductServiceID: item.ProductServiceID,
			Quantity:         item.Quantity,
			TotalPrice:       0, // akan dihitung di service
			Notes:            item.Notes,
		})
	}

	return order, items, nil
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now() // fallback, tapi bisa diubah sesuai kebutuhan
	}
	return t
}

func ToOrderResponse(order *Order, customer *customer.Customer, items []*OrderItem) *dto.OrderResponse {
	respItems := make([]dto.OrderItemResponse, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, dto.OrderItemResponse{
			ID:               item.ID,
			ProductServiceID: item.ProductServiceID,
			Quantity:         item.Quantity,
			TotalPrice:       item.TotalPrice,
			Notes:            item.Notes,
		})
	}

	return &dto.OrderResponse{
		ID:              order.ID,
		InvoiceNumber:   order.InvoiceNumber,
		StoreID:         order.StoreID,
		CustomerID:      order.CustomerID,
		CustomerName:    customer.Name,
		CustomerPhone:   customer.Phone,
		CustomerAddress: customer.Address,
		Status:          order.Status,
		Discount:        order.Discount,
		TotalPrice:      order.TotalPrice,
		PaidAmount:      order.PaidAmount,
		Change:          order.Change,
		PickupDate:      order.PickupDate.Format(time.RFC3339),
		CreatedAt:       order.CreatedAt.Format(time.RFC3339),
		CreatedBy:       order.CreatedBy,
		OrderItems:      respItems,
	}
}

func UpdateOrderModel(order *Order, req *dto.OrderRequest, updatedBy string) *Order {
	order.CustomerID = req.CustomerID
	order.Status = req.Status
	order.Discount = req.Discount
	// order.TotalPrice = req.TotalPrice
	order.PaidAmount = req.PaidAmount
	// order.Change = req.Change
	// order.PickupDate = req.PickupDate
	order.UpdatedAt = time.Now()
	order.UpdatedBy = updatedBy

	return order
}

func ToOrderItemModelsForUpdate(orderID string, reqItems []dto.OrderItemRequest) []*OrderItem {
	items := make([]*OrderItem, 0, len(reqItems))
	for _, item := range reqItems {
		items = append(items, &OrderItem{
			ID:               uuid.New().String(), // ID baru karena hapus-tambah ulang
			OrderID:          orderID,
			ProductServiceID: item.ProductServiceID,
			Quantity:         item.Quantity,
			// TotalPrice:       item.TotalPrice,
			Notes: item.Notes,
		})
	}
	return items
}
