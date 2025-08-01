package store

import (
	"context"
	"log"
	"sumunar-pos-core/internal/store/dto"
	"sumunar-pos-core/internal/userstore"
	"sumunar-pos-core/middleware"
	"sumunar-pos-core/pkg/db"
	"time"

	"github.com/jackc/pgx/v5"
)

type StoreService interface {
	Create(ctx context.Context, req *dto.StoreRequest) (*Store, error)
	FindByID(ctx context.Context, id string) (*Store, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Store, int, error)
	Update(ctx context.Context, id string, req *dto.StoreRequest) (*Store, error)
	Delete(ctx context.Context, id string) error
	UpdateLogo(ctx context.Context, storeID string, logoPath string) error
	CreateTx(ctx context.Context, req *dto.StoreRequest, userID string) (*dto.StoreResponse, error)
}

type service struct {
	repo         StoreRepository
	userStoreSvc userstore.Service
	db           db.DBTX
}

func NewService(repo StoreRepository, userStoreSvc userstore.Service, db db.DBTX) StoreService {
	return &service{repo, userStoreSvc, db}
}

func (s *service) Create(ctx context.Context, req *dto.StoreRequest) (*Store, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	store := ToStoreModel(req, userID)

	if err := s.repo.Create(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *service) CreateTx(ctx context.Context, req *dto.StoreRequest, userID string) (*dto.StoreResponse, error) {
	tx, ok := s.db.(interface {
		Begin(context.Context) (pgx.Tx, error)
	})
	if !ok {
		log.Println("❌ db does not support transactions")
		return nil, ErrUnsupportedTx
	}

	pgxTx, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer pgxTx.Rollback(ctx)

	store := ToStoreModel(req, userID)

	if err := s.repo.CreateTx(ctx, pgxTx, store); err != nil {
		return nil, err
	}

	if err := s.userStoreSvc.AssignUserToStoreTx(ctx, pgxTx, userID, store.ID, userID); err != nil {
		return nil, err
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return nil, err
	}

	return &dto.StoreResponse{
		ID:      store.ID,
		Name:    store.Name,
		Address: store.Address,
		Phone:   store.Phone,
	}, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*Store, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindAll(ctx context.Context, limit, offset int) ([]*Store, int, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *service) Update(ctx context.Context, id string, req *dto.StoreRequest) (*Store, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	store, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	store.Name = req.Name
	store.Address = req.Address
	store.Phone = req.Phone
	store.UpdatedAt = time.Now()
	store.UpdatedBy = userID

	return store, s.repo.Update(ctx, store)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) UpdateLogo(ctx context.Context, storeID string, logoPath string) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	store, err := s.repo.FindByID(ctx, storeID)
	if err != nil {
		return err
	}

	store.Logo = &logoPath
	store.UpdatedAt = time.Now()
	store.UpdatedBy = userID

	return s.repo.Update(ctx, store)
}
