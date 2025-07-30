package productservice

import (
	"sumunar-pos-core/internal/base"
	"sumunar-pos-core/internal/productservice/dto"

	"github.com/google/uuid"
)

func ToProductServiceModel(req *dto.ProductServiceRequest, createdBy string) *ProductService {
	return &ProductService{
		ID:            uuid.New().String(),
		ProductID:     req.ProductID,
		ServiceTypeID: req.ServiceTypeID,
		Unit:          req.Unit,
		Price:         req.Price,
		BaseModel: base.BaseModel{
			CreatedBy: createdBy,
		},
	}
}

func ToProductServiceResponse(service *ProductService) *dto.ProductServiceResponse {
	return &dto.ProductServiceResponse{
		ID:            service.ID,
		ProductID:     service.ProductID,
		ServiceTypeID: service.ServiceTypeID,
		Unit:          service.Unit,
		Price:         service.Price,
	}
}

func ToProductServiceListResponse(services []*ProductService) []*dto.ProductServiceResponse {
	res := make([]*dto.ProductServiceResponse, 0)
	for _, s := range services {
		res = append(res, ToProductServiceResponse(s))
	}
	return res
}
