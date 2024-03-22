package common

type ResourceDTO[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"totalCount"`
}
