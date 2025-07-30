package product

import (
	"sumunar-pos-core/internal/base"
)

type Product struct {
	ID      string `db:"id"`
	Name    string `db:"name"`
	StoreID string `db:"store_id"`
	base.BaseModel
}
