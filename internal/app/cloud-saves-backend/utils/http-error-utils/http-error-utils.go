package http_error

import "github.com/gin-gonic/gin"

func HTTPError(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{"message": message})
}