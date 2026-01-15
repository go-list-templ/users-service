package repo

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/domain/event"
)

type (
	UserRepo interface {
		Store(context.Context, entity.User) error
		All(context.Context) ([]entity.User, error)
	}

	UserAvatarRepo interface {
		Set(entity.User) entity.User
	}

	OutboxRepo interface {
		Publish(context.Context, event.Event) error
	}
)
