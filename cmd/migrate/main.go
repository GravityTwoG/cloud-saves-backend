package main

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/initializers"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	"log"
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

	roleUser := models.Role{Name: "USER"}
	roleAdmin := models.Role{Name: "ADMIN"}

	err = initializers.DB.Create(&roleUser).Error
	if err != nil {
		log.Fatal(err)
	}
	err = initializers.DB.Create(&roleAdmin).Error
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations complete")
}
