package migrator

import (
	"embed"

	"github.com/go-list-templ/users-service/internal/adapter/persistence/postgres"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Up(postgres *postgres.Postgres) error {
	sqlDB := stdlib.OpenDBFromPool(postgres.Pool)

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return err
	}

	return sqlDB.Close()
}
