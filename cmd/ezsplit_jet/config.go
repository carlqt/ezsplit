package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DatabaseURL string
	DBHost      string
	DBName      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBSchema    string
	JWTSecret   []byte
}

// NewConfig creates a new config using the environment variables.
// The variables are loaded from the .env file.
// The EnvConfig struct ensures helps with type safety and auto completion.
func NewConfig() *EnvConfig {
	dbSchema := "public"

	envDBPort := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(envDBPort)
	if err != nil {
		log.Printf("failed to convert port (%s) to int\n", envDBPort)
		panic(err)
	}

	return &EnvConfig{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		DBHost:      os.Getenv("DB_HOST"),
		DBName:      os.Getenv("DB_NAME"),
		DBPort:      dbPort,
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBSchema:    dbSchema,
		JWTSecret:   []byte(os.Getenv("JWT_SECRET")),
	}
}

func InitializeEnvVariables() {
	// During test mode, the tests aren't looking for .env in the root of the project but relative to where
	// the tests are run (./internal)

	slog.Debug("initializing environment variables")

	if os.Getenv("GO_ENV") == "test" {
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			fmt.Fprintf(os.Stderr, "Unable to identify current directory (needed to load .env)")
			os.Exit(1)
		}
		basepath := filepath.Dir(file)
		err := godotenv.Load(filepath.Join(basepath, "../../.env"))
		if err != nil {
			log.Fatal("can't load environments: ", err)
		}

		return
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}
