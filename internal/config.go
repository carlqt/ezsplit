package internal

import (
	"os"
)

type EnvConfig struct {
	DatabaseURL string
	DBHost      string
	DBName      string
	DBPort      string
	DBUser      string
	DBPassword  string
	Port        string
	JWTSecret   []byte
}

// NewConfig creates a new config using the environment variables.
// The variables are loaded from the .env file.
// The EnvConfig struct ensures helps with type safety and auto completion.
func NewConfig() *EnvConfig {
	return &EnvConfig{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		DBHost:      os.Getenv("DB_HOST"),
		DBName:      os.Getenv("DB_NAME"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		JWTSecret:   []byte(os.Getenv("JWT_SECRET")),
	}
}
