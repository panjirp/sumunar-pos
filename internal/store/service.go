package store

import (
	"context"
	"log"
	"sumunar-pos-core/internal/store/dto"
	"sumunar-pos-core/middleware"
	"time"
)

type StoreService interface {
	Create(ctx context.Context, req *dto.StoreRequest) (*Store, error)
	FindByID(ctx context.Context, id string) (*Store, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Store, int, error)
	Update(ctx context.Context, id string, req *dto.StoreRequest) (*Store, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo StoreRepository
}

func NewService(repo StoreRepository) StoreService {
	return &service{repo}
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

func (s *service) FindByID(ctx context.Context, id string) (*Store, error) {
	store, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (s *service) FindAll(ctx context.Context, limit, offset int) ([]*Store, int, error) {
	stores, total, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return stores, total, nil
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
