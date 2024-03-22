package rest_utils

import (
	"github.com/gin-gonic/gin"
)

func DecodeJSON[T any](ctx *gin.Context) (T, error) {
	var dto T
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		return dto, err
	}
	return dto, nil
}

func DecodeURI[T any](ctx *gin.Context) (T, error) {
	var dto T
	err := ctx.ShouldBindUri(&dto)
	if err != nil {
		return dto, err
	}
	return dto, nil
}

func DecodeQuery[T any](ctx *gin.Context) (T, error) {
	var dto T
	err := ctx.ShouldBindQuery(&dto)
	if err != nil {
		return dto, err
	}
	return dto, nil
}
