package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func NewMigrator(dbURL, migrationsPath string) (*Migrator, error) {
	if dbURL == "" {
		return nil, errors.New("database URL is required")
	}
	if migrationsPath == "" {
		return nil, errors.New("migrations path is required")
	}

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return &Migrator{
		migrate: m,
	}, nil
}

func (m *Migrator) Up() error {
	err := m.migrate.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("apply migrations: %w", err)
	}

	version, dirty, _ := m.migrate.Version()
	log.Printf("Migrations applied successfully (version: %d, dirty: %v)", version, dirty)
	return nil
}

func (m *Migrator) Down() error {
	err := m.migrate.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("rollback migrations: %w", err)
	}
	return nil
}

func (m *Migrator) Steps(n int) error {
	err := m.migrate.Steps(n)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migration steps: %w", err)
	}
	return nil
}

func (m *Migrator) Force(version int) error {
	return m.migrate.Force(version)
}

func (m *Migrator) Version() (version uint, dirty bool, err error) {
	return m.migrate.Version()
}

func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return fmt.Errorf("close source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("close database: %w", dbErr)
	}
	return nil
}

func RunMigrations() error {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return fmt.Errorf("DB_URL environment variable is required")
	}

	migrator, err := NewMigrator(dbURL, "file://migrations")
	if err != nil {
		return err
	}
	defer migrator.Close()

	return migrator.Up()
}
