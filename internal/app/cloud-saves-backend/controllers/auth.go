package controllers

import (
	authDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/auth"
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	http_error_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/http-error-utils"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

func NewAuth(
	authService services.AuthService, 
) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	registerDTO := authDTOs.RegisterDTO{}
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
	loginDTO := authDTOs.LoginDTO{}
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
	session := sessions.Default(ctx)
	userResponseDTO := session.Get("user")
	if userResponseDTO == nil {
		http_error_utils.HTTPError(
			ctx, http.StatusUnauthorized, 
			"UNAUTHORIZED",
		)
		return
	}
	
	ctx.JSON(http.StatusOK, &userResponseDTO)
}

func (c *authController) getUserId(ctx *gin.Context) (uint, error) {
	session := sessions.Default(ctx)
	val := session.Get("user")

	var userResponseDTO *userDTOs.UserResponseDTO 
	var ok bool

	if userResponseDTO, ok = val.(*userDTOs.UserResponseDTO); !ok || userResponseDTO == nil {
		return 0, fmt.Errorf("UNAUTHORIZED")
	}

	return userResponseDTO.Id, nil
}

func (c *authController) ChangePassword(ctx *gin.Context) {
	changePasswordDTO := authDTOs.ChangePasswordDTO{}
	ctx.Bind(&changePasswordDTO)
	
	userId, err := c.getUserId(ctx)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusUnauthorized, 
			err.Error(),
		)
		return
	}

	err = c.authService.ChangePassword(userId, &changePasswordDTO)
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
	panic("unimplemented")
}

func (c *authController) ResetPassword(ctx *gin.Context) {
	panic("unimplemented")
}