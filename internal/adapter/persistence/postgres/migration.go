package postgres

import (
	"embed"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Migration struct {
	*Postgres

	logger *zap.Logger
}

func NewMigration(postgres *Postgres, logger *zap.Logger) *Migration {
	return &Migration{postgres, logger}
}

func (m Migration) Up() error {
	sqlDB := stdlib.OpenDBFromPool(m.Pool)

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		m.logger.Error("migration up", zap.Error(err))

		return err
	}

	return sqlDB.Close()
}
