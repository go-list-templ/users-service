package transaction

import (
	"context"

	"github.com/go-list-templ/users-service/internal/adapter/persistence/postgres"
	"go.uber.org/zap"
)

type Manager struct {
	*postgres.Postgres

	logger *zap.Logger
}

func NewManager(p *postgres.Postgres, l *zap.Logger) *Manager {
	return &Manager{p, l}
}

func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.Begin(ctx)
	if err != nil {
		return err
	}

	trCtx := setTrCtx(ctx, &tx)

	defer func() {
		if r := recover(); r != nil {
			err = tx.Rollback(trCtx)
			m.logger.Warn("trm panic", zap.Error(err), zap.Any("panic", r))
		}
	}()

	if err = fn(trCtx); err != nil {
		errTx := tx.Rollback(trCtx)
		if errTx != nil {
			m.logger.Warn("trm rollback", zap.Error(errTx))
		}

		return err
	}

	return tx.Commit(trCtx)
}
