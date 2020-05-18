package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	// Migrate configuration
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"

	"github.com/figment-networks/near-indexer/config"
)

func startMigrations(cfg *config.Config) error {
	// TEMP
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	return db.Automigrate()

	log.Println("getting current directory")
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	srcDir := filepath.Join(dir, "migrations")
	srcPath := fmt.Sprintf("file://%s", srcDir)

	log.Println("using migrations from", srcDir)
	migrations, err := migrate.New(srcPath, cfg.DatabaseURL)
	if err != nil {
		return err
	}

	log.Println("running migrations")
	return migrations.Up()
}
