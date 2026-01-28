package port

import (
	"context"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/domain/event"
)

//go:generate mockgen -source=repo.go -destination=../../../test/mocks/mocks_repo_test.go -package=mock_test

type (
	UserRepo interface {
		Store(context.Context, entity.User) error
		All(context.Context) ([]entity.User, error)
	}

	OutboxRepo interface {
		Publish(context.Context, event.Event) error
	}
)
