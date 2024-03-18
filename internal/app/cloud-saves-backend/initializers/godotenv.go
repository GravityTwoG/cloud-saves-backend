package initializers

import (
	"github.com/joho/godotenv"
)

func LoadEnvVariables() (err error) {
	err = godotenv.Load();
	return err
}