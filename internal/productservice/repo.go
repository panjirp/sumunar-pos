package productservice

import (
	"context"
	"sumunar-pos-core/pkg/db"
)

type ProductServiceRepository interface {
	Create(ctx context.Context, service *ProductService) error
	FindByID(ctx context.Context, id string) (*ProductService, error)
	FindAll(ctx context.Context, limit, offset int) ([]*ProductService, int, error)
	Update(ctx context.Context, service *ProductService) error
	Delete(ctx context.Context, id string) error
}

type productServiceRepo struct {
	db db.DBTX
}

func NewServiceRepo(db db.DBTX) ProductServiceRepository {
	return &productServiceRepo{db}
}

func (r *productServiceRepo) Create(ctx context.Context, productservice *ProductService) error {
	query := `
		INSERT INTO product_service (id, product_id, service_type_id, unit, price, is_active, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		productservice.ID,
		productservice.ProductID,
		productservice.ServiceTypeID,
		productservice.Unit,
		productservice.Price,
		productservice.IsActive,
		productservice.CreatedAt,
		productservice.CreatedBy,
	)
	return err
}

func (r *productServiceRepo) FindByID(ctx context.Context, id string) (*ProductService, error) {
	query := `SELECT id, product_id, service_type_id, unit, price, is_active, created_at, created_by, updated_at, updated_by FROM product_service WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var productservice ProductService
	err := row.Scan(
		&productservice.ID,
		&productservice.ProductID,
		&productservice.ServiceTypeID,
		&productservice.Unit,
		&productservice.Price,
		&productservice.IsActive,
		&productservice.CreatedAt,
		&productservice.CreatedBy,
		&productservice.UpdatedAt,
		&productservice.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &productservice, nil
}

func (r *productServiceRepo) FindAll(ctx context.Context, limit, offset int) ([]*ProductService, int, error) {
	query := `SELECT id, product_id, service_type_id, unit, price, is_active, created_at, created_by, updated_at, updated_by FROM product_service LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var productServices []*ProductService
	for rows.Next() {
		var s ProductService
		if err := rows.Scan(
			&s.ID,
			&s.ProductID,
			&s.ServiceTypeID,
			&s.Unit,
			&s.Price,
			&s.IsActive,
			&s.CreatedAt,
			&s.CreatedBy,
			&s.UpdatedAt,
			&s.UpdatedBy,
		); err != nil {
			return nil, 0, err
		}
		productServices = append(productServices, &s)
	}

	// total count
	var total int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM product_service`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return productServices, total, nil
}

func (r *productServiceRepo) Update(ctx context.Context, productService *ProductService) error {
	query := `
		UPDATE product_service SET product_id = $1, service_type_id = $2, unit = $3, price = $4, is_active = $5,
		updated_at = $6, updated_by = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(ctx, query,
		productService.ProductID,
		productService.ServiceTypeID,
		productService.Unit,
		productService.Price,
		productService.IsActive,
		productService.UpdatedAt,
		productService.UpdatedBy,
		productService.ID,
	)
	return err
}

func (r *productServiceRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM product_service WHERE id = $1`, id)
	return err
}
