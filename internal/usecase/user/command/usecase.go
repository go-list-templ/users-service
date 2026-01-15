package command

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/domain/event"
	"github.com/go-list-templ/grpc/internal/repo"
)

type UserUsecase struct {
	userRepo   repo.UserRepo
	outboxRepo repo.OutboxRepo
}

func NewUserUsecase(u repo.UserRepo, o repo.OutboxRepo) *UserUsecase {
	return &UserUsecase{userRepo: u, outboxRepo: o}
}

func (u *UserUsecase) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := u.userRepo.Store(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	userCreated := event.NewUserCreated(user)

	err = u.outboxRepo.Publish(ctx, userCreated.Event)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
