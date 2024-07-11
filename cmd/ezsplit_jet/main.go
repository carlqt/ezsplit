// This package wraps jet-go to generate the models using a desired folder structure (excludes db name)

package main

import (
	"log"
	"strconv"

	"github.com/carlqt/ezsplit/internal"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
)

func main() {
	internal.InitializeEnvVariables()

	genDir := ".gen"
	config := internal.NewConfig()

	dbPort, err := strconv.Atoi(config.DBPort)
	if err != nil {
		log.Printf("failed to convert port (%s) to int\n", config.DBPort)
		panic(err)
	}

	dbConnection := postgres.DBConnection{
		Host:     config.DBHost,
		Port:     dbPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		SslMode:  "disable",

		DBName:     config.DBName,
		SchemaName: "dvds",
	}

	err = postgres.Generate(
		genDir,
		dbConnection,
		template.Default(postgres2.Dialect),
	)

	if err != nil {
		log.Println("failed to generate models")
		panic(err)
	}
}
