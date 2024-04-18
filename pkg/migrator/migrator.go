package migrator

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(efs embed.FS, connString string) (uint, bool, error) {
	fs, err := iofs.New(efs, ".")
	if err != nil {
		return 0, false, fmt.Errorf("iofs new: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", fs, connString)
	if err != nil {
		return 0, false, fmt.Errorf("migrate new with source instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return 0, false, fmt.Errorf("migrate up: %w", err)
	}

	v, d, err := m.Version()
	if err != nil {
		return 0, false, fmt.Errorf("migrate version: %w", err)
	}

	return v, d, nil
}
