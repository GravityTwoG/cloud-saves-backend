package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddRedirectRoutes(router *gin.RouterGroup) {
	router.GET("/redirect", redirect)
}

func redirect(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, ctx.Query("redirect-to"))
}