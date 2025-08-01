package customer

import (
	"context"
	"sumunar-pos-core/pkg/db"
)

type CustomerRepository interface {
	Create(ctx context.Context, product *Customer) error
	FindByID(ctx context.Context, id string) (*Customer, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Customer, int, error)
	Update(ctx context.Context, product *Customer) error
	Delete(ctx context.Context, id string) error
}

type customerRepo struct {
	db db.DBTX
}

func NewCustomerRepository(db db.DBTX) CustomerRepository {
	return &customerRepo{db}
}

func (r *customerRepo) Create(ctx context.Context, customer *Customer) error {
	query := `
		INSERT INTO customers (id, name, phone, address, is_active, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		customer.ID,
		customer.Name,
		customer.Phone,
		customer.Address,
		customer.IsActive,
		customer.CreatedAt,
		customer.CreatedBy,
	)
	return err
}

func (r *customerRepo) FindByID(ctx context.Context, id string) (*Customer, error) {
	query := `SELECT id, name, phone, address, is_active, created_at, created_by, updated_at, updated_by FROM customers WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var product Customer
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Phone,
		&product.Address,
		&product.IsActive,
		&product.CreatedAt,
		&product.CreatedBy,
		&product.UpdatedAt,
		&product.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *customerRepo) FindAll(ctx context.Context, limit, offset int) ([]*Customer, int, error) {
	query := `SELECT id, name, phone, address, is_active, created_at, created_by, updated_at, updated_by FROM customers LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var s Customer
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Phone,
			&s.Address,
			&s.IsActive,
			&s.CreatedAt,
			&s.CreatedBy,
			&s.UpdatedAt,
			&s.UpdatedBy,
		); err != nil {
			return nil, 0, err
		}
		customers = append(customers, &s)
	}

	// total count
	var total int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM customers`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *customerRepo) Update(ctx context.Context, product *Customer) error {
	query := `
		UPDATE customers SET name = $1, phone = $2, address = $3, is_active = $4,
		updated_at = $5, updated_by = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query,
		product.Name,
		product.Phone,
		product.Address,
		product.IsActive,
		product.UpdatedAt,
		product.UpdatedBy,
		product.ID,
	)
	return err
}

func (r *customerRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM customers WHERE id = $1`, id)
	return err
}
