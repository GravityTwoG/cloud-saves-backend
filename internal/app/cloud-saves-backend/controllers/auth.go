package controllers

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/auth"
	"cloud-saves-backend/internal/app/cloud-saves-backend/middlewares"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	auth_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/auth"
	http_error_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/http-error-utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)


func AddAuthRoutes(router *gin.RouterGroup, authService services.AuthService) {
	authController := newAuth(authService)

	authRouter := router.Group("/auth")
	
	authRouter.POST("/registration", middlewares.Unauthorized, authController.Register)
	authRouter.POST("/login", middlewares.Unauthorized, authController.Login)
	authRouter.POST("/logout", middlewares.Auth, authController.Logout)
	authRouter.GET("/me", middlewares.Auth, authController.Me)
	authRouter.POST("/auth-change-password", middlewares.Auth, authController.ChangePassword)
	authRouter.POST("/recover-password", middlewares.Unauthorized, authController.RequestPasswordReset)
	authRouter.POST("/change-password", middlewares.Unauthorized, authController.ResetPassword)
}

type AuthController interface {
	Register(*gin.Context)

	Login(*gin.Context)

	Logout(*gin.Context)

	Me(*gin.Context)

	ChangePassword(*gin.Context)

	RequestPasswordReset(*gin.Context)

	ResetPassword(*gin.Context)
}

type authController struct {
	authService services.AuthService
}


func newAuth(
	authService services.AuthService, 
) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	registerDTO := auth.RegisterDTO{}
 	err := ctx.ShouldBindJSON(&registerDTO)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest, 
			err.Error(),
		)
		return
	}

	userResponseDTO, err := c.authService.Register(&registerDTO)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest, 
			err.Error(),
		)
		return
	}

	ctx.JSON(http.StatusCreated, userResponseDTO)
}

func (c *authController) Login(ctx *gin.Context) {
	loginDTO := auth.LoginDTO{}
	ctx.Bind(&loginDTO)

	userResponseDTO, err := c.authService.Login(&loginDTO)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusUnauthorized, 
			err.Error(),
		)
		return
	}

	session := sessions.Default(ctx)
	session.Set("user", userResponseDTO) 
	session.Save()

	ctx.JSON(http.StatusOK, userResponseDTO)
}

func (c *authController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("user")
	session.Save()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out.",
	})
}

func (c *authController) Me(ctx *gin.Context) {
	userResponseDTO, err := auth_utils.ExtractUser(ctx)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusUnauthorized, 
			"UNAUTHORIZED",
		)
		return
	}
	
	ctx.JSON(http.StatusOK, &userResponseDTO)
}

func (c *authController) ChangePassword(ctx *gin.Context) {
	changePasswordDTO := auth.ChangePasswordDTO{}
	ctx.Bind(&changePasswordDTO)
	
	user, err := auth_utils.ExtractUser(ctx)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusUnauthorized, 
			err.Error(),
		)
		return
	}

	err = c.authService.ChangePassword(user.Id, &changePasswordDTO)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest, 
			err.Error(),
		)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *authController) RequestPasswordReset(ctx *gin.Context) {
	requestPasswordResetDTO := auth.RequestPasswordResetDTO{}
	err := ctx.ShouldBindJSON(&requestPasswordResetDTO)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest, 
			err.Error(),
		)
		return
	}

	err = c.authService.RequestPasswordReset(&requestPasswordResetDTO)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest, 
			err.Error(),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (c *authController) ResetPassword(ctx *gin.Context) {
	resetPasswordDTO := auth.ResetPasswordDTO{}
	err := ctx.ShouldBindJSON(&resetPasswordDTO)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest, 
			err.Error(),
		)
		return
	}

	err = c.authService.ResetPassword(&resetPasswordDTO)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest, 
			err.Error(),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
