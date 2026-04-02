package migrator

import (
	"errors"

	"github.com/go-list-templ/users-service/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationsURL = "file:///migrations"

func Up(cfg *config.DB) error {
	m, err := migrate.New(migrationsURL, cfg.URL)
	if err != nil {
		return err
	}

	err = m.Up()

	errSource, errDb := m.Close()
	if errSource != nil || errDb != nil {
		return errors.Join(errSource, errDb)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}
