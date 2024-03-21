package main

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	"errors"
	"log"
	"os"

	"gorm.io/gorm"
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

func main() {
	err := db.AutoMigrate(&models.Role{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.PasswordRecoveryToken{})
	if err != nil {
		log.Fatal(err)
	}

	roleUser := models.Role{Name: models.RoleUser}
	roleAdmin := models.Role{Name: models.RoleAdmin}

	err = db.Create(&roleUser).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}
	err = db.Create(&roleAdmin).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}

	user, err := models.NewUser(
		"admin",
		"admin@example.com",
		"12121212",
		&roleAdmin,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Create(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}

	log.Println("Migrations complete")
}
