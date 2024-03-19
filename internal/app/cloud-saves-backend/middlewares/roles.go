package middlewares

import (
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	http_error_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/http-error-utils"
	"net/http"
	"slices"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Roles(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user := session.Get("user")
		if user == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "UNAUTHORIZED",
			})
			ctx.Abort()
			return
		}
		if len(roles) == 0 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "FORBIDDEN",
			})
			ctx.Abort()
			return
		}

		userResponseDTO, ok := user.(*userDTOs.UserResponseDTO)
		if !ok {
			http_error_utils.HTTPError(
				ctx, http.StatusInternalServerError, 
				"INTERNAL_SERVER_ERROR",
			)
			ctx.Abort()
			return
		}

		if slices.Contains(roles, userResponseDTO.Role) {
			ctx.Next()
			return
		}
		
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "FORBIDDEN",
		})
		ctx.Abort()
	}
}