package main

import (
	"cloud-saves-backend/docs"
	"cloud-saves-backend/internal/app/cloud-saves-backend/config"
	"cloud-saves-backend/internal/app/cloud-saves-backend/controllers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	email_sender "cloud-saves-backend/internal/app/cloud-saves-backend/email-sender"
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/middlewares"
	"cloud-saves-backend/internal/app/cloud-saves-backend/repositories"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

var db *gorm.DB

func init() {

	err := initializers.LoadEnvVariables()
	if err != nil {
		log.Fatal(err)
	}

	db, err = initializers.ConnectToDB(os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}
}

func createApp(database *gorm.DB, conf *config.Config) *gin.Engine {
	gob.Register(&user.UserResponseDTO{})

	app := gin.New()

	app.Use(
		gin.Logger(),
		middlewares.Recovery(recoveryHandler),
		middlewares.CORS(conf.AllowedOrigins),
		middlewares.Sessions(conf.RedisHost, conf.SessionSecret),
	)

	mailer := email_sender.NewEmailSender(
		conf.EmailSenderName,
		conf.EmailSenderAddress,
		conf.EmailSenderPassword,
		conf.EmailAuthAddress,
		conf.EmailServerAddress,
	)
	apiBaseURL := conf.APIAddress + conf.APIPrefix
	emailService := services.NewEmail(mailer, apiBaseURL)

	trManager := manager.Must(
		trmgorm.NewDefaultFactory(database),
		manager.WithSettings(trmgorm.MustSettings(
			settings.Must(
				settings.WithPropagation(trm.PropagationNested))),
		),
	)

	userRepo := repositories.NewUserRepository(database, trmgorm.DefaultCtxGetter)
	roleRepo := repositories.NewRoleRepository(database, trmgorm.DefaultCtxGetter)
	recoveryTokenRepo := repositories.NewPasswordRecoveryTokenRepository(database, trmgorm.DefaultCtxGetter)

	ctx := context.Background()
	authService := services.NewAuth(trManager, ctx, roleRepo, userRepo, recoveryTokenRepo, emailService)

	apiRouter := app.Group(conf.APIPrefix)

	controllers.AddAuthRoutes(apiRouter, authService)
	controllers.AddRedirectRoutes(apiRouter)

	docs.SwaggerInfo.BasePath = conf.APIPrefix
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return app
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": "INTERNAL_SERVER_ERROR"})
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
	conf := config.LoadConfig()
	app := createApp(db, conf)

	if app == nil {
		return
	}

	app.Run()
}
