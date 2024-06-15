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
	// repositories
	// configs
	Repositories *repository.Repository
	Config       *EnvConfig
	DB           *sql.DB
}

func init() {
	if os.Getenv("GO_ENV") == "test" {
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			fmt.Fprintf(os.Stderr, "Unable to identify current directory (needed to load .env)")
			os.Exit(1)
		}
		basepath := filepath.Dir(file)
		err := godotenv.Load(filepath.Join(basepath, "../.env"))
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func NewApp() *App {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	config := NewConfig()
	db := newDB(config)
	repositories := repository.NewRepository(db)

	return &App{Config: config, DB: db, Repositories: repositories}
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
