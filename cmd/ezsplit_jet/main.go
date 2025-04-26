// This package wraps jet-go to generate the models using a desired folder structure (excludes db name)

package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-jet/jet/v2/generator/postgres"
	_ "github.com/lib/pq"
)

func main() {
	InitializeEnvVariables()

	config := NewConfig()
	genDir := ".gen"

	jetDestinationPath := path.Join(genDir, config.DBName, config.DBSchema)
	expectedPath := path.Join(genDir, config.DBSchema)

	err := jetGenerate(config, genDir)
	if err != nil {
		panic(err)
	}

	// After successful generation
	// find the <schema> folder under .gen directory
	err = moveDir(jetDestinationPath, expectedPath)
	if err != nil {
		panic(err)
	}

	// Cleanup
	dbNameDir := path.Join(genDir, config.DBName)
	os.RemoveAll(dbNameDir) //nolint:errcheck
	fmt.Printf("Generated models in %s\n", expectedPath)
}

func jetGenerate(dbConfig *EnvConfig, destinationPath string) error {
	dbConnection := postgres.DBConnection{
		Host:     dbConfig.DBHost,
		Port:     dbConfig.DBPort,
		User:     dbConfig.DBUser,
		Password: dbConfig.DBPassword,
		SslMode:  "disable",

		DBName:     dbConfig.DBName,
		SchemaName: dbConfig.DBSchema,
	}

	err := postgres.Generate(
		destinationPath,
		dbConnection,
	)

	if err != nil {
		log.Println("failed to generate models")
		return fmt.Errorf("failed to generate postgres models: %w", err)
	}

	return nil
}

func moveDir(src string, dest string) error {
	if !dirExists(src) {
		return nil
	}

	err := os.RemoveAll(dest)
	if err != nil {
		return fmt.Errorf("failed to delete destination folder when moving: %w", err)
	}

	err = os.Rename(src, dest)
	if err != nil {
		return fmt.Errorf("failed to move %s to %s: %w", src, dest, err)
	}

	return nil
}

func dirExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}
