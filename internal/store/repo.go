package store

import (
	"context"

	"sumunar-pos-core/pkg/db"
)

type StoreRepository interface {
	Create(ctx context.Context, store *Store) error
	CreateTx(ctx context.Context, tx db.DBTX, store *Store) error
	FindByID(ctx context.Context, id string) (*Store, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Store, int, error)
	Update(ctx context.Context, store *Store) error
	Delete(ctx context.Context, id string) error
}

type storeRepo struct {
	db db.DBTX
}

func NewStoreRepository(db db.DBTX) StoreRepository {
	return &storeRepo{db}
}

func (r *storeRepo) Create(ctx context.Context, store *Store) error {
	query := `
		INSERT INTO stores (id, name, code, address, phone, logo, is_active, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		store.ID,
		store.Name,
		store.Code,
		store.Address,
		store.Phone,
		store.IsActive,
		store.CreatedAt,
		store.CreatedBy,
	)
	return err
}

func (r *storeRepo) CreateTx(ctx context.Context, tx db.DBTX, store *Store) error {
	query := `
		INSERT INTO stores (id, name, code, address, phone, logo, is_active, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $8, $9)
	`
	_, err := tx.Exec(ctx, query,
		store.ID,
		store.Name,
		store.Code,
		store.Address,
		store.Phone,
		store.Logo,
		store.IsActive,
		store.CreatedAt,
		store.CreatedBy,
	)
	return err
}

func (r *storeRepo) FindByID(ctx context.Context, id string) (*Store, error) {
	query := `SELECT id, name, code, address, phone, logo, is_active, created_at, created_by, updated_at, updated_by FROM stores WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var store Store
	err := row.Scan(
		&store.ID,
		&store.Name,
		&store.Code,
		&store.Address,
		&store.Phone,
		&store.Logo,
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
	query := `SELECT id, name, code, address, phone, logo, is_active, created_at, created_by, updated_at, updated_by FROM stores LIMIT $1 OFFSET $2`
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
			&s.Code,
			&s.Address,
			&s.Phone,
			&s.Logo,
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
		UPDATE stores SET name = $1, code = $2, address = $3, phone = $4, logo = $5, is_active = $6,
		updated_at = $7, updated_by = $8
		WHERE id = $9
	`
	_, err := r.db.Exec(ctx, query,
		store.Name,
		store.Code,
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
