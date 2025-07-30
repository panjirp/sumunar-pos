package userstore

import "sumunar-pos-core/internal/base"

type UserStore struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	StoreID string `json:"store_id"`
	base.BaseModel
}
