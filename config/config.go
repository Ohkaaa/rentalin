package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret string

	XenditSecretKey     string
	XenditCallbackToken string

	ServerPort string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using system env")
	}

	return &Config{
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBName:              getEnv("DB_NAME", "rentalin"),
		JWTSecret:           getEnv("JWT_SECRET", "secret"),
		XenditSecretKey:     getEnv("XENDIT_SECRET", ""),
		XenditCallbackToken: getEnv("XENDIT_CALLBACK_TOKEN", ""),
		ServerPort:          getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
