package dto

type StoreRequest struct {
	Name    string  `json:"name" validate:"required"`
	Code    string  `json:"code" validate:"required"`
	Address string  `json:"address" validate:"required"`
	Phone   *string `json:"phone"`
	Logo    *string `json:"logo"`
}
