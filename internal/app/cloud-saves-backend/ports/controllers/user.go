package controllers

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/common"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	user_dto "cloud-saves-backend/internal/app/cloud-saves-backend/ports/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/ports/middlewares"
	http_error_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/http-error-utils"
	rest_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/rest-utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUserRoutes(router *gin.RouterGroup, userService user.UserService) {
	userController := newUser(userService)

	authRouter := router.Group("/users")

	authRouter.GET("/", middlewares.Auth, userController.GetUsers)
}

type UserController interface {
	GetUsers(*gin.Context)
}

type userController struct {
	userService user.UserService
}

func newUser(
	userService user.UserService,
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
	ctx.JSON(http.StatusOK, user_dto.FromUsers(users))
}
