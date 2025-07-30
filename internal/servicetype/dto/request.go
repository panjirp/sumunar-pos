package dto

type ServiceRequest struct {
	Name string `json:"name" validate:"required"`
}
