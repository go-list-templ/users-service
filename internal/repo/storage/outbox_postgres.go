package storage

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/event"
	"github.com/go-list-templ/grpc/pkg/postgres"
)

type OutboxPostgres struct {
	*postgres.Postgres
}

func NewOutboxPostgres(postgres *postgres.Postgres) *OutboxPostgres {
	return &OutboxPostgres{postgres}
}

func (r *OutboxPostgres) Publish(ctx context.Context, e event.Event) error {
	query := `
		INSERT INTO outbox (message_id, message) 
		VALUES ($1, $2)
	`

	_, err := r.Exec(ctx, query,
		e.ID,
		e.Payload,
	)
	if err != nil {
		return err
	}

	return nil
}
