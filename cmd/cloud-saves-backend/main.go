package main

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()
  apiRouter := app.Group("/api/v1")

	routes.AddRoutes(apiRouter)

	app.Run(":8080")
}