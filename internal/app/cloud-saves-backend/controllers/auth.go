package controllers

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/auth"
	"cloud-saves-backend/internal/app/cloud-saves-backend/middlewares"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	auth_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/auth"
	http_error_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/http-error-utils"
	json_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/json-utils"
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

// @Tags Auth
// @Summary Register new user
// @Accept json
// @Produce json
// @Param body body auth.RegisterDTO true "RegisterDTO"
// @Success 201 {object} user.UserResponseDTO
// @Router /auth/registration [post]
func (c *authController) Register(ctx *gin.Context) {
	registerDTO, err := json_utils.Decode[auth.RegisterDTO](ctx)
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

// @Tags Auth
// @Summary Login
// @Accept json
// @Produce json
// @Param body body auth.LoginDTO true "LoginDTO"
// @Success 200 {object} user.UserResponseDTO
// @Router /auth/login [post]
func (c *authController) Login(ctx *gin.Context) {
	loginDTO, err := json_utils.Decode[auth.LoginDTO](ctx)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest,
			err.Error(),
		)
		return
	}

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

// @Tags Auth
// @Summary Logout
// @Security CookieAuth
// @Produce json
// @Success 200
// @Router /auth/logout [post]
func (c *authController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("user")
	session.Save()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out.",
	})
}

// @Tags Auth
// @Summary Get current user
// @Security CookieAuth
// @Produce json
// @Success 200 {object} user.UserResponseDTO
// @Router /auth/me [get]
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

// @Tags Auth
// @Summary Change user password
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param body body auth.ChangePasswordDTO true "ChangePasswordDTO"
// @Success 200
// @Router /auth/auth-change-password [post]
func (c *authController) ChangePassword(ctx *gin.Context) {
	changePasswordDTO, err := json_utils.Decode[auth.ChangePasswordDTO](ctx)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest,
			err.Error(),
		)
		return
	}

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

// @Tags Auth
// @Summary Request password reset
// @Accept json
// @Produce json
// @Param body body auth.RequestPasswordResetDTO true "RequestPasswordResetDTO"
// @Success 200
// @Router /auth/recover-password [post]
func (c *authController) RequestPasswordReset(ctx *gin.Context) {
	requestPasswordResetDTO, err := json_utils.Decode[auth.RequestPasswordResetDTO](ctx)
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

// @Tags Auth
// @Summary Reset password
// @Accept json
// @Produce json
// @Param body body auth.ResetPasswordDTO true "ResetPasswordDTO"
// @Success 200
// @Router /auth/reset-password [post]
func (c *authController) ResetPassword(ctx *gin.Context) {
	resetPasswordDTO, err := json_utils.Decode[auth.ResetPasswordDTO](ctx)
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
