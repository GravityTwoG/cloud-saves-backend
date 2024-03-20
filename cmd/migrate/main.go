package main

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	password_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/password-utils"
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

	roleUser := models.Role{Name: "ROLE_USER"}
	roleAdmin := models.Role{Name: "ROLE_ADMIN"}

	err = db.Create(&roleUser).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}
	err = db.Create(&roleAdmin).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}

	hashedPassword, err := password_utils.HashPassword("12121212")
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  hashedPassword,
		IsBlocked: false,
		Role:      roleAdmin,
	}

	err = db.Create(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}

	log.Println("Migrations complete")
}
