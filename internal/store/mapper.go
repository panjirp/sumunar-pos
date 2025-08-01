package store

import (
	"sumunar-pos-core/internal/base"
	"sumunar-pos-core/internal/store/dto"

	"github.com/google/uuid"
)

func ToStoreModel(req *dto.StoreRequest, createdBy string) *Store {
	return &Store{
		ID:      uuid.New().String(),
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
		Phone:   req.Phone,
		Logo:    req.Logo,
		BaseModel: base.BaseModel{
			CreatedBy: createdBy,
		},
	}
}

func ToStoreResponse(store *Store) *dto.StoreResponse {
	return &dto.StoreResponse{
		ID:      store.ID,
		Name:    store.Name,
		Code:    store.Code,
		Address: store.Address,
		Phone:   store.Phone,
		Logo:    store.Logo,
	}
}

func ToStoreListResponse(stores []*Store) []*dto.StoreResponse {
	res := make([]*dto.StoreResponse, 0)
	for _, s := range stores {
		res = append(res, ToStoreResponse(s))
	}
	return res
}
