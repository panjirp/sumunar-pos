package user

import (
	"context"
	"log"

	"sumunar-pos-core/internal/base"
	"sumunar-pos-core/middleware"

	"github.com/google/uuid"
)

type Service interface {
	GetUser(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*User, int, error)
	CreateUser(ctx context.Context, username, email string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	UpdateLastLogin(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *service) FindByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.FindByUsername(ctx, username)
}

func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]*User, int, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *service) CreateUser(ctx context.Context, username, email string) (*User, error) {

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		log.Println("failed to get user id from context:", err)
	}
	u := &User{
		ID:       uuid.New().String(),
		Username: username,
		Email:    email,
		Provider: "local",
		Role:     "owner",
		BaseModel: base.BaseModel{
			CreatedBy: userID,
		},
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) UpdateLastLogin(ctx context.Context, id string) error {
	return s.repo.UpdateLastLogin(ctx, id)
}
