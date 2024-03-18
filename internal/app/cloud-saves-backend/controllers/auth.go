package controllers

import (
	authDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/auth"
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	"cloud-saves-backend/internal/app/cloud-saves-backend/sessions"
	"net/http"
	"time"

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

func NewAuth(authService services.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	registerDTO := authDTOs.RegisterDTO{}
	ctx.Bind(&registerDTO)

	userResponseDTO, err := c.authService.Register(&registerDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, userResponseDTO)
}

func (c *authController) Login(ctx *gin.Context) {
	loginDTO := authDTOs.LoginDTO{}
	ctx.Bind(&loginDTO)

	userResponseDTO, err := c.authService.Login(&loginDTO)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	session := sessions.Create(userResponseDTO)

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   session.Id,
		Expires: session.ExpiresAt,
		Path: "/",
		Domain: "",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	ctx.Writer.Header().Add("Set-Cookie", cookie.String())
	ctx.JSON(http.StatusOK, userResponseDTO)
}

func (c *authController) Logout(ctx *gin.Context) {
	cookie, err := ctx.Request.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unauthorized",
	}	)
		return
	}
	
	_, err = sessions.Get(cookie.Value)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
	}	)
		return
	}

	sessions.Delete(cookie.Value)
	
	cookie = &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Now(),
		Path: "/",
		Domain: "",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	ctx.Writer.Header().Add("Set-Cookie", cookie.String())
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out.",
	})
}

func (c *authController) Me(ctx *gin.Context) {
	cookie, err := ctx.Request.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unauthorized",
	}	)
		return
	}
	
	session, err := sessions.Get(cookie.Value)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	userResponseDTO := userDTOs.UserResponseDTO{
		Id: session.UserId,
		Email: session.Email,
		Username: session.Username,
		Role: session.Role,
		IsBlocked: false,
	}
	
	ctx.JSON(http.StatusOK, &userResponseDTO)
}

func (c *authController) getUserId(ctx *gin.Context) (uint, error) {
	cookie, err := ctx.Request.Cookie("session_id")
	if err != nil {
		return 0, err
	}
	
	session, err := sessions.Get(cookie.Value)
	if err != nil {
		return 0, err
	}

	return session.UserId, nil
}

func (c *authController) ChangePassword(ctx *gin.Context) {
	changePasswordDTO := authDTOs.ChangePasswordDTO{}
	ctx.Bind(&changePasswordDTO)
	
	userId, err := c.getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = c.authService.ChangePassword(userId, &changePasswordDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
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