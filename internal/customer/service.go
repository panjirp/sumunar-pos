package customer

import (
	"context"
	"log"
	"sumunar-pos-core/internal/customer/dto"
	"sumunar-pos-core/middleware"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerService interface {
	Create(ctx context.Context, req *dto.CustomerRequest) (*Customer, error)
	FindByID(ctx context.Context, id string) (*Customer, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Customer, int, error)
	Update(ctx context.Context, id string, req *dto.CustomerRequest) (*Customer, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo CustomerRepository
	db   *pgxpool.Pool
}

func NewService(repo CustomerRepository, db *pgxpool.Pool) CustomerService {
	return &service{repo, db}
}

func (s *service) Create(ctx context.Context, req *dto.CustomerRequest) (*Customer, error) {

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	product := ToCustomerModel(req, userID)

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*Customer, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *service) FindAll(ctx context.Context, limit, offset int) ([]*Customer, int, error) {
	products, total, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (s *service) Update(ctx context.Context, id string, req *dto.CustomerRequest) (*Customer, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.Phone = req.Phone
	product.Address = req.Address
	product.UpdatedAt = time.Now()
	product.UpdatedBy = userID

	return product, s.repo.Update(ctx, product)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
