package infrastructure

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

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

	separator := "?"
	if strings.Contains(dbURL, "?") {
		separator = "&"
	}
	dbURL = dbURL + separator + "x-migrations-table=ai_service_schema_migrations"
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
			slog.Info("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("apply migrations: %w", err)
	}

	version, dirty, _ := m.migrate.Version()
	slog.Info("Migrations applied successfully", slog.Uint64("version", uint64(version)), slog.Bool("dirty", dirty))
	return nil
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
	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		return fmt.Errorf("POSTGRES_URL environment variable is required")
	}

	migrator, err := NewMigrator(dbURL, "file://migrations")
	if err != nil {
		return err
	}
	defer migrator.Close()

	return migrator.Up()
}
