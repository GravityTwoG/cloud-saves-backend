package main

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/infra/models"
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"errors"
	"log"
	"os"

	"gorm.io/gorm"
)

func main() {
	err := initializers.LoadEnvVariables(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initializers.ConnectToDB(os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.Role{})
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

	roleUser := models.Role{Name: user.RoleUser}
	roleAdmin := models.Role{Name: user.RoleAdmin}

	err = db.Create(&roleUser).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}
	err = db.Create(&roleAdmin).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}

	user, err := user.NewUser(
		"admin",
		"admin@example.com",
		"12121212",
		user.RoleFromDB(roleAdmin.ID, roleAdmin.Name),
	)
	userModel := models.UserFromEntity(user)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Create(&userModel).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Fatal(err)
	}

	log.Println("Migrations complete")
}
