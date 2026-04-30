package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Migrator struct {
	srcDriver source.Driver
}

func NewMigrator(sqlFiles embed.FS, dirName string) (*Migrator, error) {
	d, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		return nil, err
	}
	return &Migrator{
		srcDriver: d,
	}, nil
}

func (m *Migrator) ApplyMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver migration: %w", err)
	}
	migrator, err := migrate.NewWithInstance("migration_embed_sql_files", m.srcDriver, "psql_db", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	defer migrator.Close()
	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to up migration: %w", err)
	}
	return nil
}
