package dto

type CustomerResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
