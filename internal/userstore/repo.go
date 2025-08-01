package userstore

import (
	"context"

	"sumunar-pos-core/pkg/db"
)

type UserStoreRepository interface {
	Assign(ctx context.Context, us *UserStore) error
	CreateTx(ctx context.Context, tx db.DBTX, us *UserStore) error
	FindStoresByUserID(ctx context.Context, userID string) ([]string, error)
}

type userStoreRepo struct {
	db db.DBTX
}

func NewUserStoreRepository(db db.DBTX) UserStoreRepository {
	return &userStoreRepo{db: db}
}

func (r *userStoreRepo) Assign(ctx context.Context, us *UserStore) error {
	query := `
		INSERT INTO user_stores (user_id, store_id, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, us.UserID, us.StoreID, us.CreatedAt, us.CreatedBy, us.UpdatedAt, us.UpdatedBy)
	return err
}

func (r *userStoreRepo) CreateTx(ctx context.Context, tx db.DBTX, us *UserStore) error {
	query := `
		INSERT INTO user_stores (user_id, store_id, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.Exec(ctx, query, us.UserID, us.StoreID, us.CreatedAt, us.CreatedBy, us.UpdatedAt, us.UpdatedBy)
	return err
}

func (r *userStoreRepo) FindStoresByUserID(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT store_id FROM user_stores WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []string
	for rows.Next() {
		var storeID string
		if err := rows.Scan(&storeID); err != nil {
			return nil, err
		}
		stores = append(stores, storeID)
	}
	return stores, nil
}
