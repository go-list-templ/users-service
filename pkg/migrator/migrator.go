package migrator

import (
	"errors"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/go-list-templ/users-service/pkg/config"
	"github.com/golang-migrate/migrate/v4"
)

const migrationsUrl = "file:///migrations"

func Up(cfg *config.DB) error {
	m, err := migrate.New(migrationsUrl, cfg.URL)
	if err != nil {
		return err
	}

	err = m.Up()

	errSource, errDb := m.Close()
	if errSource != nil || errDb != nil {
		return errors.Join(errSource, errDb)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}
