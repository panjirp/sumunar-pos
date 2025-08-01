package dto

type CustomerRequest struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}
