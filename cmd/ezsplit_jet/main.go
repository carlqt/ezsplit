// This package wraps jet-go to generate the models using a desired folder structure (excludes db name)

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/carlqt/ezsplit/internal"
	"github.com/go-jet/jet/v2/generator/postgres"
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

	dbSchema := "public"

	dbConnection := postgres.DBConnection{
		Host:     config.DBHost,
		Port:     dbPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		SslMode:  "disable",

		DBName:     config.DBName,
		SchemaName: dbSchema,
	}

	err = postgres.Generate(
		genDir,
		dbConnection,
	)

	if err != nil {
		log.Println("failed to generate models")
		panic(err)
	}

	// After successful generation
	// find the <schema> folder under .gen directory
	srcDir := path.Join(genDir, config.DBName, dbSchema)
	destDir := path.Join(genDir, dbSchema)

	err = moveDir(srcDir, destDir)
	if err != nil {
		panic(err)
	}

	dbNameDir := path.Join(genDir, config.DBName)
	os.RemoveAll(dbNameDir)
	fmt.Printf("Generated models in %s\n", destDir)
}

func moveDir(src string, dest string) error {
	err := os.RemoveAll(dest)
	if err != nil {
		return fmt.Errorf("failed to delete destination folder when moving: %w", err)
	}

	err = os.Rename(src, dest)
	if err != nil {
		return fmt.Errorf("failed to move %s to %s\n: %w", src, dest, err)
	}

	return nil
}
