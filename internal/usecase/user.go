package usecase

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo"
)

type User struct {
	repo       repo.UserRepo
	avatarRepo repo.UserAvatarRepo
}

func NewUser(r repo.UserRepo, a repo.UserAvatarRepo) *User {
	return &User{repo: r, avatarRepo: a}
}

func (u *User) All(ctx context.Context) ([]entity.User, error) {
	users, err := u.repo.All(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := u.repo.Store(ctx, user)
	if err != nil {
		return user, err
	}

	u.avatarRepo.Set(user)

	return user, nil
}
