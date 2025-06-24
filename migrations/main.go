package migrations

import (
	"database/sql"
	"golang-tutorial/config"
)

type migration interface {
	Name() string
	Up(conn *sql.Tx) error
	Down(conn *sql.Tx) error
	SkipProd() bool
}

func getMigrations() []migration {
	return []migration{
		getCreatePerpusTable(),
		getCreateDownloadTokensTable(),
	}
}

func checkDuplicateMigrationNames(migrations []migration) {
	nameSet := make(map[string]bool)
	for _, m := range migrations {
		if nameSet[m.Name()] {
			panic("Duplicate migrations name found: " + m.Name())
		}
		nameSet[m.Name()] = true
	}
}

func Up(db *sql.DB) {
	// Get all migrations
	migrations := getMigrations()

	cfg := config.Get()

	// Check for duplicate migrations names
	checkDuplicateMigrationNames(migrations)

	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			name VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)

	if err != nil {
		panic(err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	for _, m := range migrations {
		// Check if migrations has already been applied
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", m.Name()).Scan(&count)
		if err != nil {
			panic(err)
		}

		if count == 0 {
			if cfg.IsProduction && m.SkipProd() {
				continue
			}

			// Apply migrations
			if err := m.Up(tx); err != nil {
				panic(err)
			}

			// Record migrations as applied
			_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", m.Name())
			if err != nil {
				panic(err)
			}

			println("Applied migrations:", m.Name())
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func Down(db *sql.DB) {
	// Get all migrations
	migrations := getMigrations()

	// Check for duplicate migrations names
	checkDuplicateMigrationNames(migrations)

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	// Get the last applied migrations
	var lastMigration string
	err = tx.QueryRow("SELECT name FROM migrations ORDER BY applied_at DESC LIMIT 1").Scan(&lastMigration)
	if err != nil {
		if err == sql.ErrNoRows {
			println("No migrations to revert")
			return
		}
		panic(err)
	}

	// Find the migrations in our list
	var migrationToRevert migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if migrations[i].Name() == lastMigration {
			migrationToRevert = migrations[i]
			break
		}
	}

	if migrationToRevert == nil {
		panic("Last applied migrations not found in migrations list")
	}

	// Revert the migrations
	if err := migrationToRevert.Down(tx); err != nil {
		panic(err)
	}

	// Remove the migrations record
	_, err = tx.Exec("DELETE FROM migrations WHERE name = $1", lastMigration)
	if err != nil {
		panic(err)
	}

	println("Reverted migrations:", lastMigration)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func DownAll(db *sql.DB) {
	// Get all migrations
	migrations := getMigrations()

	// Check for duplicate migrations names
	checkDuplicateMigrationNames(migrations)

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		// Check if migrations has been applied
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", m.Name()).Scan(&count)
		if err != nil {
			panic(err)
		}

		if count > 0 {
			// Revert migrations
			if err := m.Down(tx); err != nil {
				panic(err)
			}

			// Remove migrations record
			_, err = tx.Exec("DELETE FROM migrations WHERE name = $1", m.Name())
			if err != nil {
				panic(err)
			}

			println("Reverted migrations:", m.Name())
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
