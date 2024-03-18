package routes

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/modules/users"

	"github.com/gin-gonic/gin"
)

func usersRoutes(router *gin.RouterGroup) {
	usersRouter := router.Group("/users")
	
	usersRouter.GET("/me", users.Me)
}