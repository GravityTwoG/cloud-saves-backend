package controllers

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/common"
	"cloud-saves-backend/internal/app/cloud-saves-backend/middlewares"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	http_error_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/http-error-utils"
	rest_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/rest-utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUserRoutes(router *gin.RouterGroup, userService services.UserService) {
	userController := newUser(userService)

	authRouter := router.Group("/users")

	authRouter.GET("/", middlewares.Auth, userController.GetUsers)
}

type UserController interface {
	GetUsers(*gin.Context)
}

type userController struct {
	userService services.UserService
}

func newUser(
	userService services.UserService,
) UserController {
	return &userController{
		userService: userService,
	}
}

// @Tags Users
// @Summary Get users
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param pageNumber query int true "Page number"
// @Param pageSize query int true "Page size"
// @Param searchQuery query string false "Search query"
// @Success 200 {object} user.UsersResponseDTO
// @Router /users [get]
func (c *userController) GetUsers(ctx *gin.Context) {
	dto, err := rest_utils.DecodeQuery[common.GetResourceDTO](ctx)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	users, err := c.userService.GetUsers(ctx, dto)
	if err != nil {
		http_error_utils.HTTPError(
			ctx, http.StatusInternalServerError,
			err.Error(),
		)
		return
	}
	ctx.JSON(http.StatusOK, users)
}
