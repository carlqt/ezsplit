package internal

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/carlqt/ezsplit/internal/repository"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type App struct {
	Repositories *repository.Repository
	Config       *EnvConfig
	DB           *sql.DB
}

// InitializeEnvVariables uses godot export the Env variables from .env files
// In production, this could be optimized by presuming the env variables are already available and avoid calling this function
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
		err := godotenv.Load(filepath.Join(basepath, "../.env"))
		if err != nil {
			log.Fatal("can't load environments", err)
		}

		return
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func InitializeLogger() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

// Initializes Config, DB and Repositories that will be used by the server
func NewApp() *App {
	InitializeLogger()
	InitializeEnvVariables()

	config := NewConfig()
	db := newDB(config)
	repositories := repository.NewRepository(db)

	return &App{
		Config:       config,
		DB:           db,
		Repositories: repositories,
	}
}

func newDB(config *EnvConfig) *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBName, config.DBPassword)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
