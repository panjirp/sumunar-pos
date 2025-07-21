package store

import (
	"sumunar-pos-core/internal/base"
)

type Store struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Phone   *string `json:"phone,omitempty"`
	base.BaseModel
}
