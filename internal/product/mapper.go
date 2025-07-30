package product

import (
	"sumunar-pos-core/internal/base"
	"sumunar-pos-core/internal/product/dto"

	"github.com/google/uuid"
)

func ToProductModel(req *dto.ProductRequest, createdBy string) *Product {
	return &Product{
		ID:      uuid.New().String(),
		Name:    req.Name,
		StoreID: req.StoreID,
		BaseModel: base.BaseModel{
			CreatedBy: createdBy,
		},
	}
}

func ToProductResponse(product *Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:      product.ID,
		Name:    product.Name,
		StoreID: product.StoreID,
	}
}

func ToProductListResponse(products []*Product) []*dto.ProductResponse {
	res := make([]*dto.ProductResponse, 0)
	for _, s := range products {
		res = append(res, ToProductResponse(s))
	}
	return res
}
