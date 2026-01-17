package command

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/domain/event"
	"github.com/go-list-templ/grpc/internal/repo"
	"github.com/go-list-templ/grpc/pkg/trm"
)

type UserUsecase struct {
	userRepo   repo.UserRepo
	outboxRepo repo.OutboxRepo

	trm *trm.Manager
}

func NewUserUsecase(u repo.UserRepo, o repo.OutboxRepo, t *trm.Manager) *UserUsecase {
	return &UserUsecase{u, o, t}
}

func (u *UserUsecase) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := u.trm.Do(ctx, func(ctx context.Context) error {
		err := u.userRepo.Store(ctx, user)
		if err != nil {
			return err
		}

		userCreated := event.NewUserCreated(user)

		return u.outboxRepo.Publish(ctx, userCreated.Event)
	})
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
