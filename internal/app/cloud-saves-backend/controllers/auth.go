package controllers

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/auth"
	"cloud-saves-backend/internal/app/cloud-saves-backend/middlewares"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	auth_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/auth"
	http_error_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/http-error-utils"
	rest_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/rest-utils"
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

	onlyAdmin := middlewares.Roles([]models.RoleName{models.RoleAdmin})

	authRouter.POST("/block-user/:userId", onlyAdmin, authController.BlockUser)
	authRouter.POST("/unblock-user/:userId", onlyAdmin, authController.UnblockUser)
}

type AuthController interface {
	Register(*gin.Context)

	Login(*gin.Context)

	Logout(*gin.Context)

	Me(*gin.Context)

	ChangePassword(*gin.Context)

	RequestPasswordReset(*gin.Context)

	ResetPassword(*gin.Context)

	BlockUser(*gin.Context)

	UnblockUser(*gin.Context)
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
	registerDTO, err := rest_utils.DecodeJSON[auth.RegisterDTO](ctx)
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
	loginDTO, err := rest_utils.DecodeJSON[auth.LoginDTO](ctx)
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
	changePasswordDTO, err := rest_utils.DecodeJSON[auth.ChangePasswordDTO](ctx)
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
	requestPasswordResetDTO, err := rest_utils.DecodeJSON[auth.RequestPasswordResetDTO](ctx)
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
	resetPasswordDTO, err := rest_utils.DecodeJSON[auth.ResetPasswordDTO](ctx)
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

// @Tags Auth
// @Summary Block user
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param userId path auth.BlockUserDTO true "User ID"
// @Success 200
// @Router /auth/block-user/{userId} [post]
func (c *authController) BlockUser(ctx *gin.Context) {
	dto, err := rest_utils.DecodeURI[auth.BlockUserDTO](ctx)
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
	if user.Id == dto.UserId {
		http_error_utils.HTTPError(
			ctx, http.StatusForbidden,
			"CANNOT_BLOCK_YOURSELF",
		)
		return
	}

	err = c.authService.BlockUser(dto.UserId)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	// TODO: clear all sessions of user

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

// @Tags Auth
// @Summary Unblock user
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param userId path auth.BlockUserDTO true "BlockUserDTO"
// @Success 200
// @Router /auth/unblock-user/{userId} [post]
func (c *authController) UnblockUser(ctx *gin.Context) {
	dto, err := rest_utils.DecodeURI[auth.BlockUserDTO](ctx)
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
			ctx, http.StatusForbidden,
			err.Error(),
		)
		return
	}
	if user.Id == dto.UserId {
		http_error_utils.HTTPError(
			ctx, http.StatusForbidden,
			"CANNOT_UNBLOCK_YOURSELF",
		)
		return
	}

	err = c.authService.UnblockUser(dto.UserId)
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
