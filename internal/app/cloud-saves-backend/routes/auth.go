package routes

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/controllers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/middlewares"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"

	"github.com/gin-gonic/gin"
)

func authRoutes(router *gin.RouterGroup) {
	authService := services.NewAuth(initializers.DB)
	authController := controllers.NewAuth(authService)
	authMiddleware := middlewares.Auth()
	unauthorizedMiddleware := middlewares.Unauthorized()

	authRouter := router.Group("/auth")

	authRouter.POST("/registration", unauthorizedMiddleware, authController.Register)
	authRouter.POST("/login", unauthorizedMiddleware, authController.Login)
	authRouter.POST("/logout", authMiddleware, authController.Logout)
	authRouter.GET("/me", authMiddleware, authController.Me)
	authRouter.POST("/auth-change-password", authMiddleware, authController.ChangePassword)
	authRouter.POST("/recover-password", unauthorizedMiddleware, authController.RequestPasswordReset)
	authRouter.POST("/change-password", unauthorizedMiddleware, authController.ResetPassword)
}