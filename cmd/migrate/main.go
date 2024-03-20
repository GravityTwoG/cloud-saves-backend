package main

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	password_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/password-utils"
	"errors"
	"log"

	"gorm.io/gorm"
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
	err := initializers.DB.AutoMigrate(&models.Role{})
	if err != nil {
		log.Fatal(err)
	}

	err = initializers.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}
	
	err = initializers.DB.AutoMigrate(&models.PasswordRecoveryToken{})
	if err != nil {
		log.Fatal(err)
	}

	roleUser := models.Role{Name: "ROLE_USER"}
	roleAdmin := models.Role{Name: "ROLE_ADMIN"}

	err = initializers.DB.Create(&roleUser).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}
	err = initializers.DB.Create(&roleAdmin).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}

	hashedPassword, err := password_utils.HashPassword("12121212");
	if err != nil {
		log.Fatal(err)
	}
	
	user := models.User{
		Username: "admin",
		Email: "admin@example.com",
		Password: hashedPassword,
		IsBlocked: false,
		Role: roleAdmin,
	}
	
	err = initializers.DB.Create(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}

	log.Println("Migrations complete")
}
