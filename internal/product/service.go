package product

import (
	"context"
	"log"
	"sumunar-pos-core/internal/product/dto"
	"sumunar-pos-core/middleware"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService interface {
	Create(ctx context.Context, req *dto.ProductRequest) (*Product, error)
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Product, int, error)
	Update(ctx context.Context, id string, req *dto.ProductRequest) (*Product, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo ProductRepository
	db   *pgxpool.Pool
}

func NewService(repo ProductRepository, db *pgxpool.Pool) ProductService {
	return &service{repo, db}
}

func (s *service) Create(ctx context.Context, req *dto.ProductRequest) (*Product, error) {

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	product := ToProductModel(req, userID)

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *service) FindAll(ctx context.Context, limit, offset int) ([]*Product, int, error) {
	products, total, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (s *service) Update(ctx context.Context, id string, req *dto.ProductRequest) (*Product, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.StoreID = req.StoreID
	product.UpdatedAt = time.Now()
	product.UpdatedBy = userID

	return product, s.repo.Update(ctx, product)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
