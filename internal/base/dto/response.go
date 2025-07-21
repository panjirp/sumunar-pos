package dto

type PaginatedResponse[T any] struct {
	Data   []T `json:"data"`
	Total  int `json:"total"`
	Page   int `json:"page"`
	Pages  int `json:"pages"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
