package tests

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"log"
	"os"

	"gorm.io/gorm"
)

func SetupSuite() *gorm.DB {
	err := initializers.LoadEnvVariables("../../../../../.env.test")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initializers.ConnectToDB(os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
