package initializers

import (
	"github.com/joho/godotenv"
)

func LoadEnvVariables(filename string) (err error) {
	err = godotenv.Load(filename)
	return err
}
