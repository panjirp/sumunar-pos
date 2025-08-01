package servicetype

import (
	"context"
	"log"
	"sumunar-pos-core/internal/servicetype/dto"
	"sumunar-pos-core/middleware"
	"sumunar-pos-core/pkg/db"
	"time"
)

type ServiceTypeService interface {
	Create(ctx context.Context, req *dto.ServiceRequest) (*ServiceType, error)
	FindByID(ctx context.Context, id string) (*ServiceType, error)
	FindAll(ctx context.Context, limit, offset int) ([]*ServiceType, int, error)
	Update(ctx context.Context, id string, req *dto.ServiceRequest) (*ServiceType, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo ServiceRepository
	db   db.DBTX
}

func NewService(repo ServiceRepository, db db.DBTX) ServiceTypeService {
	return &service{repo, db}
}

func (s *service) Create(ctx context.Context, req *dto.ServiceRequest) (*ServiceType, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	service := ToServiceTypeModel(req, userID)

	if err := s.repo.Create(ctx, service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*ServiceType, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindAll(ctx context.Context, limit, offset int) ([]*ServiceType, int, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *service) Update(ctx context.Context, id string, req *dto.ServiceRequest) (*ServiceType, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	service, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	service.Name = req.Name
	service.UpdatedAt = time.Now()
	service.UpdatedBy = userID

	return service, s.repo.Update(ctx, service)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
