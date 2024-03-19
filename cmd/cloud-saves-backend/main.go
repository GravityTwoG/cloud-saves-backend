package main

import (
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/middlewares"
	"cloud-saves-backend/internal/app/cloud-saves-backend/routes"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
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
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(middlewares.Recovery(recoveryHandler))
	
	config := cors.DefaultConfig()
  config.AllowOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	config.AllowCredentials = true
	app.Use(cors.New(config))

	store, err := redis.NewStore(10, "tcp", os.Getenv("REDIS_HOST"), "", []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		log.Fatal(err)
		return
	}
	store.Options(sessions.Options{
		HttpOnly: true, 
		MaxAge: 86400, 
		SameSite: http.SameSiteNoneMode,
		Secure: true,
		Path: "/",
	})
	gob.Register(&userDTOs.UserResponseDTO{})
  app.Use(sessions.Sessions("session", store))

  apiRouter := app.Group("")

	routes.AddRoutes(apiRouter)

	app.Run()
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": "INTERNAL_SERVER_ERROR"})
}