package port

import (
	"context"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/domain/event"
)

type (
	UserRepo interface {
		Store(context.Context, entity.User) error
		All(context.Context) ([]entity.User, error)
	}

	OutboxRepo interface {
		Publish(context.Context, event.Event) error
	}
)
