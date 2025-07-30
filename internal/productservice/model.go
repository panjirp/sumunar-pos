package productservice

import (
	"sumunar-pos-core/internal/base"
)

type ProductService struct {
	ID            string  `db:"id"`
	ProductID     string  `db:"product_id"`
	ServiceTypeID string  `db:"service_type_id"`
	Unit          string  `db:"unit"`  // "kg", "pcs", "m2", dll
	Price         float64 `db:"price"` // harga per unit
	base.BaseModel
}
