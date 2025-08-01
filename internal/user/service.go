package user

import (
	"context"
	"log"

	"sumunar-pos-core/internal/user/dto"
	"sumunar-pos-core/middleware"
)

type Service interface {
	GetUser(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*User, int, error)
	CreateUser(ctx context.Context, req *dto.UserRequest) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByFullname(ctx context.Context, fullname string) (*User, error)
	UpdateLastLogin(ctx context.Context, id string) error
}

type service struct {
	repo UserRepository
}

func NewService(repo UserRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *service) FindByFullname(ctx context.Context, fullname string) (*User, error) {
	return s.repo.FindByFullname(ctx, fullname)
}

func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]*User, int, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *service) CreateUser(ctx context.Context, req *dto.UserRequest) (*User, error) {

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}

	u := ToUserModel(req, userID)

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) UpdateLastLogin(ctx context.Context, id string) error {
	return s.repo.UpdateLastLogin(ctx, id)
}
