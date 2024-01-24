package internal

import (
	"os"
)

type EnvConfig struct {
	DatabaseURL string
	Port        string
}

// NewConfig creates a new config using the environment variables.
// The variables are loaded from the .env file.
// The EnvConfig struct ensures helps with type safety and auto completion.
func NewConfig() *EnvConfig {
	return &EnvConfig{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}
}
