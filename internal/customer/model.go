package customer

import (
	"sumunar-pos-core/internal/base"
)

type Customer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	base.BaseModel
}
