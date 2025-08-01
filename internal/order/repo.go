package order

import (
	"context"

	"sumunar-pos-core/pkg/db"
)

type OrderRepository interface {
	Create(ctx context.Context, tx db.DBTX, order *Order, items []*OrderItem) error
	FindByID(ctx context.Context, id string) (*Order, []*OrderItem, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Order, int, error)
	Update(ctx context.Context, tx db.DBTX, order *Order, items []*OrderItem) error
	Delete(ctx context.Context, id string) error
	CountTodayOrders(ctx context.Context, storeID string) (int, error)
	FindOrderItemsByOrderIDs(ctx context.Context, orderIDs []string) (map[string][]*OrderItem, error)
}

type orderRepo struct {
	db db.DBTX
}

func NewOrderRepository(db db.DBTX) OrderRepository {
	return &orderRepo{db}
}

func (r *orderRepo) Create(ctx context.Context, tx db.DBTX, order *Order, items []*OrderItem) error {
	query := `
		INSERT INTO orders (id, store_id, invoice_number, customer_id, status, discount, total_price, paid_amount, change, pickup_date, created_at, created_by, updated_at, updated_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$11,$12)
	`
	_, err := tx.Exec(ctx, query,
		order.ID,
		order.StoreID,
		order.InvoiceNumber,
		order.CustomerID,
		order.Status,
		order.Discount,
		order.TotalPrice,
		order.PaidAmount,
		order.Change,
		order.PickupDate,
		order.CreatedAt,
		order.CreatedBy,
	)
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err := tx.Exec(ctx, `
			INSERT INTO order_items (id, order_id, product_service_id, quantity, total_price, notes)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			item.ID,
			item.OrderID,
			item.ProductServiceID,
			item.Quantity,
			item.TotalPrice,
			item.Notes,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *orderRepo) FindByID(ctx context.Context, id string) (*Order, []*OrderItem, error) {
	query := `
		SELECT id, store_id, invoice_number, customer_id, status, discount, total_price, paid_amount, change, pickup_date, created_at, created_by, updated_at, updated_by
		FROM orders
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var o Order
	err := row.Scan(
		&o.ID,
		&o.StoreID,
		&o.InvoiceNumber,
		&o.CustomerID,
		&o.Status,
		&o.Discount,
		&o.TotalPrice,
		&o.PaidAmount,
		&o.Change,
		&o.PickupDate,
		&o.CreatedAt,
		&o.CreatedBy,
		&o.UpdatedAt,
		&o.UpdatedBy,
	)
	if err != nil {
		return nil, nil, err
	}

	items := []*OrderItem{}
	rows, err := r.db.Query(ctx, `SELECT id, order_id, product_service_id, quantity, total_price, notes FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.ProductServiceID,
			&i.Quantity,
			&i.TotalPrice,
			&i.Notes,
		); err != nil {
			return nil, nil, err
		}
		items = append(items, &i)
	}

	return &o, items, nil
}

func (r *orderRepo) FindAll(ctx context.Context, limit, offset int) ([]*Order, int, error) {
	query := `
		SELECT id, store_id, invoice_number, customer_id, status, discount, total_price, paid_amount, change, pickup_date, created_at, created_by, updated_at, updated_by
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.ID,
			&o.StoreID,
			&o.InvoiceNumber,
			&o.CustomerID,
			&o.Status,
			&o.Discount,
			&o.TotalPrice,
			&o.PaidAmount,
			&o.Change,
			&o.PickupDate,
			&o.CreatedAt,
			&o.CreatedBy,
			&o.UpdatedAt,
			&o.UpdatedBy,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, &o)
	}

	var total int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepo) Update(ctx context.Context, tx db.DBTX, order *Order, items []*OrderItem) error {
	_, err := tx.Exec(ctx, `
		UPDATE orders SET
			store_id = $1,
			invoice_number = $2,
			customer_id = $3,
			status = $4,
			discount = $5,
			total_price = $6,
			paid_amount = $7,
			change = $8,
			pickup_date = $9,
			updated_at = $10,
			updated_by = $11
		WHERE id = $12
	`,
		order.StoreID,
		order.InvoiceNumber,
		order.CustomerID,
		order.Status,
		order.Discount,
		order.TotalPrice,
		order.PaidAmount,
		order.Change,
		order.PickupDate,
		order.UpdatedAt,
		order.UpdatedBy,
		order.ID,
	)
	if err != nil {
		return err
	}

	// Delete existing items
	_, err = tx.Exec(ctx, `DELETE FROM order_items WHERE order_id = $1`, order.ID)
	if err != nil {
		return err
	}

	// Insert updated items
	for _, item := range items {
		_, err := tx.Exec(ctx, `
			INSERT INTO order_items (id, order_id, product_service_id, quantity, total_price, notes)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			item.ID,
			item.OrderID,
			item.ProductServiceID,
			item.Quantity,
			item.TotalPrice,
			item.Notes,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *orderRepo) Delete(ctx context.Context, id string) error {
	// delete order items first
	_, err := r.db.Exec(ctx, `DELETE FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return err
	}
	// delete order
	_, err = r.db.Exec(ctx, `DELETE FROM orders WHERE id = $1`, id)
	return err
}

func (r *orderRepo) CountTodayOrders(ctx context.Context, storeID string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM orders 
		WHERE store_id = $1 AND DATE(created_at) = CURRENT_DATE
	`
	var count int
	err := r.db.QueryRow(ctx, query, storeID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *orderRepo) FindOrderItemsByOrderIDs(ctx context.Context, orderIDs []string) (map[string][]*OrderItem, error) {
	// Buat query IN
	query := `SELECT id, order_id, product_service_id, quantity, total_price, notes FROM order_items WHERE order_id = ANY($1)`
	rows, err := r.db.Query(ctx, query, orderIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]*OrderItem)
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductServiceID,
			&item.Quantity,
			&item.TotalPrice,
			&item.Notes,
		); err != nil {
			return nil, err
		}
		result[item.OrderID] = append(result[item.OrderID], &item)
	}

	return result, nil
}
