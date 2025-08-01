package order

import (
	"context"
	"fmt"
	"time"

	"sumunar-pos-core/internal/customer"
	"sumunar-pos-core/internal/order/dto"
	"sumunar-pos-core/internal/productservice"

	"sumunar-pos-core/pkg/db"
)

type OrderService struct {
	repo               OrderRepository
	productServiceRepo productservice.ProductServiceRepository
	customerRepo       customer.CustomerRepository
	db                 db.TxBeginner
}

func NewService(repo OrderRepository, productServiceRepo productservice.ProductServiceRepository, customerRepo customer.CustomerRepository, db db.TxBeginner) *OrderService {
	return &OrderService{repo, productServiceRepo, customerRepo, db}
}

func (s *OrderService) CreateOrder(ctx context.Context, dto dto.OrderRequest, createdBy string) (*dto.OrderResponse, error) {
	order, items, err := ToOrderModel(dto, createdBy)
	if err != nil {
		return nil, err
	}

	cust, err := s.customerRepo.FindByID(ctx, order.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Generate invoice number
	order.InvoiceNumber, err = s.GenerateInvoiceNumber(ctx, order.StoreID)
	if err != nil {
		return nil, err
	}

	// Hitung total harga berdasarkan product_service_id dan quantity
	var total float64
	for _, item := range items {
		ps, err := s.productServiceRepo.FindByID(ctx, item.ProductServiceID)
		if err != nil {
			return nil, fmt.Errorf("product service not found: %w", err)
		}
		item.TotalPrice = ps.Price * item.Quantity
		total += item.TotalPrice
	}

	// Hitung diskon
	discountAmount := total * (order.Discount / 100)
	order.TotalPrice = total - discountAmount

	// Hitung kembalian
	order.Change = order.PaidAmount - order.TotalPrice

	// Simpan ke DB dalam transaksi
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // auto rollback jika ada error

	if err := s.repo.Create(ctx, tx, order, items); err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	res := ToOrderResponse(order, cust, items)
	return res, nil
}

func (s *OrderService) GenerateInvoiceNumber(ctx context.Context, storeID string) (string, error) {
	today := time.Now().Format("060102") // YYMMDD
	count, err := s.repo.CountTodayOrders(ctx, storeID)
	if err != nil {
		return "", err
	}
	// e.g. INV-240730-001
	return fmt.Sprintf("INV-%s-%03d", today, count+1), nil
}

func (s *OrderService) FindAll(ctx context.Context, limit, offset int) ([]*dto.OrderResponse, int, error) {
	orders, total, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Ambil semua order_id
	orderIDs := make([]string, len(orders))
	for i, o := range orders {
		orderIDs[i] = o.ID
	}

	// Ambil semua item sekaligus
	itemsMap, err := s.repo.FindOrderItemsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		cust, _ := s.customerRepo.FindByID(ctx, order.CustomerID)
		items := itemsMap[order.ID]
		res := ToOrderResponse(order, cust, items)
		responses = append(responses, res)
	}

	return responses, total, nil
}

func (s *OrderService) Update(ctx context.Context, id string, req *dto.OrderRequest, updatedBy string) (*dto.OrderResponse, error) {
	// Ambil order lama
	order, orderItems, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	cust, err := s.customerRepo.FindByID(ctx, order.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Hitung ulang total & change
	var total float64
	for _, item := range req.Items {
		total += item.TotalPrice
	}
	total = total - (total * req.Discount / 100)
	change := req.PaidAmount - total

	// Parse pickup date
	pickupDate, err := time.Parse(time.RFC3339, req.PickupDate)
	if err != nil {
		return nil, fmt.Errorf("invalid pickup_date: %w", err)
	}

	// Update order model
	UpdateOrderModel(order, req, updatedBy)
	order.TotalPrice = total
	order.Change = change
	order.PickupDate = pickupDate

	// Buat item model baru
	orderItems = ToOrderItemModelsForUpdate(order.ID, req.Items)

	// Begin TX
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.repo.Update(ctx, tx, order, orderItems); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return ToOrderResponse(order, cust, orderItems), nil
}

func (s *OrderService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
