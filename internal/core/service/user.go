package service

import (
	"context"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/event"
	"github.com/go-list-templ/users-service/internal/core/dto"
	"github.com/go-list-templ/users-service/internal/port"
	"github.com/go-list-templ/users-service/pkg/paginate"
)

type User struct {
	userRepo   port.UserRepo
	outboxRepo port.OutboxRepo
	trm        port.TransactionManager
}

func NewUser(u port.UserRepo, o port.OutboxRepo, t port.TransactionManager) *User {
	return &User{u, o, t}
}

func (s *User) Create(ctx context.Context, input dto.CreateInput) (entity.User, error) {
	user, err := entity.NewUser(input.Name, input.Email)
	if err != nil {
		return entity.User{}, err
	}

	err = s.trm.Do(ctx, func(ctx context.Context) error {
		err = s.userRepo.Store(ctx, user)
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

func (s *User) List(ctx context.Context, input dto.ListInput) (dto.ListOutput, error) {
	pagination := paginate.NewUUIDPaginate(input.PageToken)

	return s.userRepo.List(ctx, pagination)
}
