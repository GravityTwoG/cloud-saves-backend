package middlewares

import (
	auth_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	user, err := auth_utils.ExtractUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "UNAUTHORIZED",
		})
		ctx.Abort()
		return
	}

	if user.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "USER_IS_BLOCKED",
		})
		ctx.Abort()
		return
	}
	
	ctx.Next()
}