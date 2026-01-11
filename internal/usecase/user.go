package usecase

import (
	"context"
	"github.com/go-list-templ/grpc/internal/domain/event"
	"github.com/go-list-templ/grpc/pkg/uow"
	"github.com/jackc/pgx/v5"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo"
)

type User struct {
	userRepo   repo.UserRepo
	outboxRepo repo.OutboxRepo
	uow        *uow.UnitOfWork
}

func NewUser(u repo.UserRepo, o repo.OutboxRepo, uo *uow.UnitOfWork) *User {
	return &User{userRepo: u, outboxRepo: o, uow: uo}
}

func (u *User) All(ctx context.Context) ([]entity.User, error) {
	users, err := u.userRepo.All(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := u.uow.Do(ctx, func(tx pgx.Tx) error {
		err := u.userRepo.Store(ctx, tx, user)
		if err != nil {
			return err
		}

		userCreated := event.NewUserCreated(user)

		return u.outboxRepo.Publish(ctx, tx, userCreated.Event)
	})
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
