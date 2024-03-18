package routes

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/controllers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"

	"github.com/gin-gonic/gin"
)

func authRoutes(router *gin.RouterGroup) {
	authService := services.NewAuth(initializers.DB)
	authController := controllers.NewAuth(authService)

	authRouter := router.Group("/auth")
	
	authRouter.POST("/registration", authController.Register)
	authRouter.POST("/login", authController.Login)
	authRouter.POST("/logout", authController.Logout)
	authRouter.GET("/me", authController.Me)
	authRouter.POST("/auth-change-password", authController.ChangePassword)
	authRouter.POST("/recover-password", authController.RequestPasswordReset)
	authRouter.POST("/change-password", authController.ChangePassword)
}