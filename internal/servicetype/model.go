package servicetype

import (
	"sumunar-pos-core/internal/base"
)

type ServiceType struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	IsActive bool   `db:"is_active"`
	base.BaseModel
}
