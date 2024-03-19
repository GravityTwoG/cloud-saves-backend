package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Unauthorized() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user := session.Get("user")
		if user != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "FORBIDDEN",
			})
			ctx.Abort()
			return
		}
		
		ctx.Next()
	}
}