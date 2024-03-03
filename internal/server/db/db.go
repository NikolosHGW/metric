package db

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	err = runMigrations(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func runMigrations(db *sqlx.DB) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	instance, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get instance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs:migrations", d, "postgres", instance)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}

	return nil
}
