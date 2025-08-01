package productservice

import (
	"context"
	"log"
	"sumunar-pos-core/internal/productservice/dto"
	"sumunar-pos-core/middleware"
	"sumunar-pos-core/pkg/db"
	"time"
)

type ProductServiceService interface {
	Create(ctx context.Context, req *dto.ProductServiceRequest) (*ProductService, error)
	FindByID(ctx context.Context, id string) (*ProductService, error)
	FindAll(ctx context.Context, limit, offset int) ([]*ProductService, int, error)
	Update(ctx context.Context, id string, req *dto.ProductServiceRequest) (*ProductService, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo ProductServiceRepository
	db   db.DBTX
}

func NewService(repo ProductServiceRepository, db db.DBTX) ProductServiceService {
	return &service{repo: repo, db: db}
}

func (s *service) Create(ctx context.Context, req *dto.ProductServiceRequest) (*ProductService, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	productService := ToProductServiceModel(req, userID)

	if err := s.repo.Create(ctx, productService); err != nil {
		return nil, err
	}

	return productService, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*ProductService, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindAll(ctx context.Context, limit, offset int) ([]*ProductService, int, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *service) Update(ctx context.Context, id string, req *dto.ProductServiceRequest) (*ProductService, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	productService, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	productService.ProductID = req.ProductID
	productService.ServiceTypeID = req.ServiceTypeID
	productService.Unit = req.Unit
	productService.Price = req.Price
	productService.UpdatedAt = time.Now()
	productService.UpdatedBy = userID

	return productService, s.repo.Update(ctx, productService)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
