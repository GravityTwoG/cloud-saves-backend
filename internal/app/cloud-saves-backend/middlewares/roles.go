package middlewares

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/utils/auth"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func Roles(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := auth.ExtractUser(ctx)
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

		if !slices.Contains(roles, user.Role) {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "FORBIDDEN",
			})
			ctx.Abort()
			return
		}
		
		ctx.Next()
	}
}