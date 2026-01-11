package repo

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/domain/event"
	"github.com/jackc/pgx/v5"
)

type (
	UserRepo interface {
		Store(context.Context, pgx.Tx, entity.User) error
		All(context.Context) ([]entity.User, error)
	}

	UserAvatarRepo interface {
		Set(entity.User) entity.User
	}

	OutboxRepo interface {
		Publish(context.Context, pgx.Tx, event.Event) error
	}
)
