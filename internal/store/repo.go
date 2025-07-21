package store

import (
	"context"

	"sumunar-pos-core/pkg/db"
)

type StoreRepository interface {
	Create(ctx context.Context, store *Store) error
	FindByID(ctx context.Context, id string) (*Store, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Store, int, error)
	Update(ctx context.Context, store *Store) error
	Delete(ctx context.Context, id string) error
}

type storeRepo struct {
	db db.DBTX
}

func NewStoreRepo(db db.DBTX) StoreRepository {
	return &storeRepo{db}
}

func (r *storeRepo) Create(ctx context.Context, store *Store) error {
	query := `
		INSERT INTO stores (id, name, address, phone, is_active, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		store.ID,
		store.Name,
		store.Address,
		store.Phone,
		store.IsActive,
		store.CreatedAt,
		store.CreatedBy,
	)
	return err
}

func (r *storeRepo) FindByID(ctx context.Context, id string) (*Store, error) {
	query := `SELECT id, name, address, phone, is_active, created_at, created_by, updated_at, updated_by FROM stores WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var store Store
	err := row.Scan(
		&store.ID,
		&store.Name,
		&store.Address,
		&store.Phone,
		&store.IsActive,
		&store.CreatedAt,
		&store.CreatedBy,
		&store.UpdatedAt,
		&store.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepo) FindAll(ctx context.Context, limit, offset int) ([]*Store, int, error) {
	query := `SELECT id, name, address, phone, is_active, created_at, created_by, updated_at, updated_by FROM stores LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stores []*Store
	for rows.Next() {
		var s Store
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Address,
			&s.Phone,
			&s.IsActive,
			&s.CreatedAt,
			&s.CreatedBy,
			&s.UpdatedAt,
			&s.UpdatedBy,
		); err != nil {
			return nil, 0, err
		}
		stores = append(stores, &s)
	}

	// total count
	var total int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM stores`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

func (r *storeRepo) Update(ctx context.Context, store *Store) error {
	query := `
		UPDATE stores SET name = $1, address = $2, phone = $3, is_active = $4,
		updated_at = $5, updated_by = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query,
		store.Name,
		store.Address,
		store.Phone,
		store.IsActive,
		store.UpdatedAt,
		store.UpdatedBy,
		store.ID,
	)
	return err
}

func (r *storeRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM stores WHERE id = $1`, id)
	return err
}
