package repo

import (
	"context"

	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres/transaction"
	"github.com/go-list-templ/grpc/internal/core/domain/event"
	"github.com/go-list-templ/grpc/pkg/postgres"
)

type OutboxRepo struct {
	*postgres.Postgres

	getter *transaction.TrmGetter
}

func NewOutboxRepo(p *postgres.Postgres, g *transaction.TrmGetter) *OutboxRepo {
	return &OutboxRepo{p, g}
}

func (r *OutboxRepo) Publish(ctx context.Context, e event.Event) error {
	query := `INSERT INTO outbox (message_id, message) VALUES ($1, $2)`

	_, err := r.getter.TrOrDB(ctx, r.Postgres).Exec(ctx, query, e.ID, e.Payload)

	return err
}
