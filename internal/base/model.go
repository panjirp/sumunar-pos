package base

import "time"

type BaseModel struct {
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}
