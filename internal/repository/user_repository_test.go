package repository

import (
	"database/sql"
	"testing"

	"github.com/carlqt/ezsplit/config"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateWithAccount(t *testing.T) {
	config := config.NewConfig()
	db := NewDB(
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBName,
		config.DBPassword,
		"disable",
	)
	userRepository := UserRepository{db}

	t.Run("when arguments are empty strings", func(t *testing.T) {
		defer truncateAllTables(db, t)

		var pgError *pq.Error

		_, err := userRepository.CreateWithAccount("", "")

		if assert.ErrorAs(t, err, &pgError) {
			assert.Equal(t, pgError.Constraint, "non_empty_username")
		}
	})

	t.Run("when username already exists", func(t *testing.T) {
		defer truncateAllTables(db, t)

		var pgError *pq.Error

		userRepository.CreateWithAccount("john_doe", "password")
		_, err := userRepository.CreateWithAccount("john_doe", "password")

		if assert.ErrorAs(t, err, &pgError) {
			assert.Equal(t, pgError.Constraint, "idx_accounts_on_username")
		}
	})
}

// importing the truncateAllTables from integration_test package causes an import cycle
func truncateAllTables(db *sql.DB, t *testing.T) {
	query := `
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname =current_schema()) LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`

	_, err := db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
