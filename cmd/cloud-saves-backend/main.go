package main

import (
	"cloud-saves-backend/docs"
	email_sender "cloud-saves-backend/internal/app/cloud-saves-backend/adapters/email-sender"
	"cloud-saves-backend/internal/app/cloud-saves-backend/adapters/repositories"
	"cloud-saves-backend/internal/app/cloud-saves-backend/adapters/services"
	sessions_store "cloud-saves-backend/internal/app/cloud-saves-backend/adapters/sessions"
	"cloud-saves-backend/internal/app/cloud-saves-backend/config"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/auth"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/ports/controllers"
	user_dto "cloud-saves-backend/internal/app/cloud-saves-backend/ports/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/ports/middlewares"
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

func createApp(database *gorm.DB, conf *config.Config) *gin.Engine {
	gob.Register(&user_dto.UserDTO{})

	app := gin.New()

	sessionsStore, err := sessions_store.NewStore(10, "tcp", conf.RedisHost, "", conf.SessionSecret)
	if err != nil {
		log.Fatal(err)
	}
	app.Use(
		gin.Logger(),
		middlewares.Recovery(recoveryHandler),
		middlewares.CORS(conf.AllowedOrigins),
		middlewares.Sessions(sessionsStore),
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

	authService := auth.NewAuth(
		trManager,
		roleRepo,
		userRepo,
		recoveryTokenRepo,
		emailService,
	)
	userService := user.NewUserService(userRepo, roleRepo)

	apiRouter := app.Group(conf.APIPrefix)

	controllers.AddAuthRoutes(apiRouter, authService, sessionsStore)
	controllers.AddUserRoutes(apiRouter, userService)
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
	err := initializers.LoadEnvVariables(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initializers.ConnectToDB(os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	conf := config.LoadConfig()
	app := createApp(db, conf)

	if app == nil {
		return
	}

	app.Run()
}
