package internal

import (
	"database/sql"
	"fmt"
	"log"

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

func NewApp() *App {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	config := NewConfig()
	db := newDB(config)
	repositories := repository.NewRepository(db)

	return &App{Config: config, DB: db, Repositories: repositories}
}

// TODO: Split the databsae config into host, dbname, port, user and password
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
