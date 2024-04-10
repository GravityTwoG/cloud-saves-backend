package common

// ResourceDTO[T]
// @Description: generic resource DTO
type ResourceDTO[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"totalCount"`
}
