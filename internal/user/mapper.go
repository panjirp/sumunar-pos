package user

import (
	"time"

	"sumunar-pos-core/internal/base"
	basedto "sumunar-pos-core/internal/base/dto"
	"sumunar-pos-core/internal/user/dto"

	"github.com/google/uuid"
)

func ToUserModel(req *dto.UserRequest, createdBy string) *User {
	return &User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		GoogleID: req.GoogleID,
		Picture:  req.Picture,
		Provider: req.Provider,
		Role:     req.Role,
		BaseModel: base.BaseModel{
			CreatedAt: time.Now(),
			CreatedBy: createdBy,
		},
	}
}

func NewPaginatedResponse[T any](data []T, total, page, limit int) basedto.PaginatedResponse[T] {
	pages := 0
	if limit > 0 {
		pages = (total + limit - 1) / limit // ceil
	}

	offset := (page - 1) * limit
	return basedto.PaginatedResponse[T]{
		Data:   data,
		Total:  total,
		Page:   page,
		Pages:  pages,
		Limit:  limit,
		Offset: offset,
	}
}
