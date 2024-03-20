package middlewares

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/utils/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Unauthorized(ctx *gin.Context) {
	_, err := auth.ExtractUser(ctx)
	if err == nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "FORBIDDEN",
		})
		ctx.Abort()
		return
	}
	
	ctx.Next()
}