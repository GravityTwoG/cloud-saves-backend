package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddRedirectRoutes(router *gin.RouterGroup) {
	redirectController := newRedirect()

	router.GET("/redirect", redirectController.Redirect)
}

type RedirectController interface {
	Redirect(*gin.Context)
}

type redirectController struct{}

func newRedirect() RedirectController {
	return &redirectController{}
}

// @Tags Redirect
// @Summary Redirect to a given URL
// @Param   redirect-to query string true "Redirect URL"
// @Success 302 {string} string "Redirected"
// @Router /redirect [get]
func (r *redirectController) Redirect(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, ctx.Query("redirect-to"))
}
