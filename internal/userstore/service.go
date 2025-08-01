package userstore

import (
	"context"

	"sumunar-pos-core/internal/base"
	"sumunar-pos-core/pkg/db"

	"github.com/google/uuid"
)

type Service interface {
	AssignUserToStore(ctx context.Context, userID, storeID, actorID string) error
	AssignUserToStoreTx(ctx context.Context, tx db.DBTX, userID, storeID, createdBy string) error
	GetUserStoreIDs(ctx context.Context, userID string) ([]string, error)
}

type service struct {
	repo UserStoreRepository
}

func NewService(repo UserStoreRepository) Service {
	return &service{repo: repo}
}

func (s *service) AssignUserToStore(ctx context.Context, userID, storeID, actorID string) error {
	us := &UserStore{
		ID:      uuid.New().String(),
		UserID:  userID,
		StoreID: storeID,
		BaseModel: base.BaseModel{
			CreatedBy: actorID,
			UpdatedBy: actorID,
		},
	}
	return s.repo.Assign(ctx, us)
}

func (s *service) AssignUserToStoreTx(ctx context.Context, tx db.DBTX, userID, storeID, createdBy string) error {
	us := &UserStore{
		ID:        uuid.New().String(),
		UserID:    userID,
		StoreID:   storeID,
		BaseModel: base.BaseModel{CreatedBy: createdBy},
	}
	return s.repo.CreateTx(ctx, tx, us)
}

func (s *service) GetUserStoreIDs(ctx context.Context, userID string) ([]string, error) {
	return s.repo.FindStoresByUserID(ctx, userID)
}
