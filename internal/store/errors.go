package store

import "errors"

var (
	ErrStoreNotFound = errors.New("store not found")
	ErrUnsupportedTx = errors.New("database does not support transactions")
)
