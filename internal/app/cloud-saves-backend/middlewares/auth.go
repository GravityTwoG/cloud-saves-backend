package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
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
		
		ctx.Next()
	}
}