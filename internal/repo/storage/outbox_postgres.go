package storage

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/event"
	"github.com/go-list-templ/grpc/pkg/postgres"
	"github.com/go-list-templ/grpc/pkg/trm"
)

type OutboxPostgres struct {
	*postgres.Postgres

	getter *trm.CtxGetter
}

func NewOutboxPostgres(p *postgres.Postgres, g *trm.CtxGetter) *OutboxPostgres {
	return &OutboxPostgres{p, g}
}

func (r *OutboxPostgres) Publish(ctx context.Context, e event.Event) error {
	query := `INSERT INTO outbox (message_id, message) VALUES ($1, $2)`

	_, err := r.getter.TrOrDB(ctx, r.Postgres).
		Exec(ctx, query,
			e.ID,
			e.Payload,
		)

	return err
}
