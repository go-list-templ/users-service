package storage

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/event"
	"github.com/go-list-templ/grpc/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type OutboxPostgres struct {
	*postgres.Postgres
}

func NewOutboxPostgres(postgres *postgres.Postgres) *OutboxPostgres {
	return &OutboxPostgres{postgres}
}

func (r *OutboxPostgres) Publish(ctx context.Context, tx pgx.Tx, e event.Event) error {
	query := `
		INSERT INTO outbox (message_id, message) 
		VALUES ($1, $2)
	`

	_, err := tx.Exec(ctx, query,
		e.ID,
		e.Payload,
	)
	if err != nil {
		return err
	}

	return nil
}
