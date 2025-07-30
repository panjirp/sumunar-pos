package servicetype

import (
	"sumunar-pos-core/internal/base"
	"sumunar-pos-core/internal/servicetype/dto"

	"github.com/google/uuid"
)

func ToServiceTypeModel(req *dto.ServiceRequest, createdBy string) *ServiceType {
	return &ServiceType{
		ID:   uuid.New().String(),
		Name: req.Name,
		BaseModel: base.BaseModel{
			CreatedBy: createdBy,
		},
	}
}

func ToServiceTypeResponse(service *ServiceType) *dto.ServiceResponse {
	return &dto.ServiceResponse{
		ID:   service.ID,
		Name: service.Name,
	}
}

func ToServiceTypeListResponse(services []*ServiceType) []*dto.ServiceResponse {
	res := make([]*dto.ServiceResponse, 0)
	for _, s := range services {
		res = append(res, ToServiceTypeResponse(s))
	}
	return res
}
