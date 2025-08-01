package product

import (
	"context"
	"sumunar-pos-core/pkg/db"
)

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Product, int, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}

type productRepo struct {
	db db.DBTX
}

func NewProductRepository(db db.DBTX) ProductRepository {
	return &productRepo{db}
}

func (r *productRepo) Create(ctx context.Context, productType *Product) error {
	query := `
		INSERT INTO products (id, name, store_id, is_active, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		productType.ID,
		productType.Name,
		productType.StoreID,
		productType.IsActive,
		productType.CreatedAt,
		productType.CreatedBy,
	)
	return err
}

func (r *productRepo) FindByID(ctx context.Context, id string) (*Product, error) {
	query := `SELECT id, name, store_id, is_active, created_at, created_by, updated_at, updated_by FROM products WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var product Product
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.StoreID,
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

func (r *productRepo) FindAll(ctx context.Context, limit, offset int) ([]*Product, int, error) {
	query := `SELECT id, name, store_id, is_active, created_at, created_by, updated_at, updated_by FROM products LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var s Product
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.StoreID,
			&s.IsActive,
			&s.CreatedAt,
			&s.CreatedBy,
			&s.UpdatedAt,
			&s.UpdatedBy,
		); err != nil {
			return nil, 0, err
		}
		products = append(products, &s)
	}

	// total count
	var total int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepo) Update(ctx context.Context, product *Product) error {
	query := `
		UPDATE products SET name = $1, store_id = $2, is_active = $3,
		updated_at = $4, updated_by = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query,
		product.Name,
		product.StoreID,
		product.IsActive,
		product.UpdatedAt,
		product.UpdatedBy,
		product.ID,
	)
	return err
}

func (r *productRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}
