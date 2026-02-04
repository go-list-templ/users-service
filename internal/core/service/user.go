package service

import (
	"context"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/domain/event"
	"github.com/go-list-templ/grpc/internal/port"
)

type User struct {
	userRepo   port.UserRepo
	outboxRepo port.OutboxRepo
	trm        port.TransactionManager
}

func NewUser(u port.UserRepo, o port.OutboxRepo, t port.TransactionManager) *User {
	return &User{u, o, t}
}

func (s *User) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := s.trm.Do(ctx, func(ctx context.Context) error {
		err := s.userRepo.Store(ctx, user)
		if err != nil {
			return err
		}

		userCreated, err := event.NewUserCreated(user)
		if err != nil {
			return err
		}

		return s.outboxRepo.Publish(ctx, userCreated.Event)
	})
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *User) All(ctx context.Context) ([]entity.User, error) {
	users, err := s.userRepo.All(ctx)
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}
