package customer

import (
	"sumunar-pos-core/internal/base"
	"sumunar-pos-core/internal/customer/dto"

	"github.com/google/uuid"
)

func ToCustomerModel(req *dto.CustomerRequest, createdBy string) *Customer {
	return &Customer{
		ID:      uuid.New().String(),
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
		BaseModel: base.BaseModel{
			CreatedBy: createdBy,
		},
	}
}

func ToCustomerResponse(customer *Customer) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		ID:      customer.ID,
		Name:    customer.Name,
		Phone:   customer.Phone,
		Address: customer.Address,
	}
}

func ToCustomerListResponse(customers []*Customer) []*dto.CustomerResponse {
	res := make([]*dto.CustomerResponse, 0)
	for _, s := range customers {
		res = append(res, ToCustomerResponse(s))
	}
	return res
}
