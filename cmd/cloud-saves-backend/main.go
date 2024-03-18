package main

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	err := initializers.LoadEnvVariables()
	if err != nil {
		log.Fatal(err)
	}

	err = initializers.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := gin.Default()
  apiRouter := app.Group("")

	routes.AddRoutes(apiRouter)

	app.Run()
}