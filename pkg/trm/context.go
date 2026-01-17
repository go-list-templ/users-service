package trm

import (
	"context"
	"github.com/go-list-templ/grpc/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type CtxKey struct{}

func setTrCtx(ctx context.Context, tx *pgx.Tx) context.Context {
	return context.WithValue(ctx, CtxKey{}, tx)
}

type CtxGetter struct {
	trm *Manager
}

func NewCtxGetter(trm *Manager) *CtxGetter {
	return &CtxGetter{trm}
}

func (c *CtxGetter) TrOrDB(ctx context.Context, db *postgres.Postgres) Tr {
	if tx, ok := ctx.Value(CtxKey{}).(*pgx.Tx); ok {
		return *tx
	}

	return db
}
