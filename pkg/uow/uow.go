package uow

import (
	"context"
	"errors"
	"github.com/go-list-templ/grpc/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var ErrTx = errors.New("failed to execute transaction")

type UnitOfWork struct {
	*postgres.Postgres

	logger *zap.Logger
}

func NewUnitOfWork(p *postgres.Postgres, l *zap.Logger) *UnitOfWork {
	return &UnitOfWork{p, l}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := u.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			u.logger.Panic("recovered from panic")
		}
	}()

	if err = fn(tx); err != nil {
		_ = tx.Rollback(ctx)

		return ErrTx
	}

	return tx.Commit(ctx)
}
