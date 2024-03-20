package main

import (
	"cloud-saves-backend/docs"
	"cloud-saves-backend/internal/app/cloud-saves-backend/controllers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	email_sender "cloud-saves-backend/internal/app/cloud-saves-backend/email-sender"
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/middlewares"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	gob.Register(&user.UserResponseDTO{})
}

// @title           Cloud Saves API
// @version         1.0
// @description     This is a cloud saves backend API

// @contact.name   Marsel Abazbekov
// @contact.url    https://github.com/GravityTwoG
// @contact.email  marsel.ave@gmail.com

// @host      localhost:8080
// @BasePath  /
// @securitydefinitions.apikey CookieAuth
// @in cookie
// @name session 
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
	
  app.Use(sessions.Sessions("session", store))

	apiPrefix := os.Getenv("API_PREFIX")
	host := os.Getenv("API_ADDRESS")
	apiBaseURL := host + apiPrefix

	docs.SwaggerInfo.BasePath = apiPrefix
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
  apiRouter := app.Group(apiPrefix)
	
	mailer := email_sender.NewEmailSender(
		os.Getenv("EMAIL_SENDER_NAME"), 
		os.Getenv("EMAIL_SENDER_ADDRESS"), 
		os.Getenv("EMAIL_SENDER_PASSWORD"), 
		os.Getenv("EMAIL_AUTH_ADDRESS"), 
		os.Getenv("EMAIL_SERVER_ADDRESS"),
	)
	emailService := services.NewEmail(mailer, apiBaseURL)
	authService := services.NewAuth(initializers.DB, emailService)
	
	controllers.AddAuthRoutes(apiRouter, authService)
	controllers.AddRedirectRoutes(apiRouter)
	
	app.Run()
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": "INTERNAL_SERVER_ERROR"})
}