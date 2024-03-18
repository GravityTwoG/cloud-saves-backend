package routes

import "github.com/gin-gonic/gin"

func AddRoutes(router *gin.RouterGroup) {
	authRoutes(router)
}