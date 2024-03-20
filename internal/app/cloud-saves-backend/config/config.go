package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                int
	APIAddress          string
	APIPrefix           string
	DSN                 string
	RedisHost           string
	AllowedOrigins      []string
	SessionSecret       []byte
	EmailSenderName     string
	EmailSenderAddress  string
	EmailSenderPassword string
	EmailAuthAddress    string
	EmailServerAddress  string
}

func LoadConfig() *Config {
	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)

	if err != nil {
		log.Println("Using default port: 8080")
		port = 8080
	}

	return &Config{
		Port:                int(port),
		APIAddress:          os.Getenv("API_ADDRESS"),
		APIPrefix:           os.Getenv("API_PREFIX"),
		DSN:                 os.Getenv("DSN"),
		RedisHost:           os.Getenv("REDIS_HOST"),
		AllowedOrigins:      strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		SessionSecret:       []byte(os.Getenv("SESSION_SECRET")),
		EmailSenderName:     os.Getenv("EMAIL_SENDER_NAME"),
		EmailSenderAddress:  os.Getenv("EMAIL_SENDER_ADDRESS"),
		EmailSenderPassword: os.Getenv("EMAIL_SENDER_PASSWORD"),
		EmailAuthAddress:    os.Getenv("EMAIL_AUTH_ADDRESS"),
		EmailServerAddress:  os.Getenv("EMAIL_SERVER_ADDRESS"),
	}
}
