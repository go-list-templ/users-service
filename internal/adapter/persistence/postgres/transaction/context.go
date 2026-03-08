package transaction

import (
	"context"

	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres"
	"github.com/jackc/pgx/v5"
)

type CtxKey struct{}

func setTrCtx(ctx context.Context, tx *pgx.Tx) context.Context {
	return context.WithValue(ctx, CtxKey{}, tx)
}

type TrmGetter struct {
	trm *Manager
}

func NewTrmGetter(trm *Manager) *TrmGetter {
	return &TrmGetter{trm}
}

func (c *TrmGetter) TrOrDB(ctx context.Context, db *postgres.Postgres) Tr {
	if tx, ok := ctx.Value(CtxKey{}).(*pgx.Tx); ok {
		return *tx
	}

	return db.Master
}
