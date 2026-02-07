package repo

import (
	"context"
	"time"

	"github.com/go-list-templ/grpc/internal/adapter/cache/redis"
	"github.com/go-list-templ/grpc/internal/adapter/cache/redis/repo/dao"
	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/port"
	"go.uber.org/zap"
)

const (
	KeyAllUsers           = "users:all"
	DefaultContextTimeout = 5 * time.Second
)

type UserRepo struct {
	repo   port.UserRepo
	redis  *redis.Redis
	logger *zap.Logger
}

func NewUserRepo(repo port.UserRepo, redis *redis.Redis, logger *zap.Logger) *UserRepo {
	return &UserRepo{repo: repo, redis: redis, logger: logger}
}

func (u *UserRepo) All(ctx context.Context) ([]entity.User, error) {
	var cachedUsers []dao.User

	err := u.redis.GetCache(ctx, KeyAllUsers, &cachedUsers)
	if err == nil && len(cachedUsers) > 0 {
		users := make([]entity.User, len(cachedUsers))

		for i, user := range cachedUsers {
			users[i] = user.ToEntity()
		}

		return users, nil
	}

	users, err := u.repo.All(ctx)
	if err != nil {
		return nil, err
	}

	go u.cacheAllUsers(users)

	return users, nil
}

func (u *UserRepo) cacheAllUsers(users []entity.User) {
	defer func() {
		if r := recover(); r != nil {
			u.logger.Error("panic in cacheAllUsers", zap.Any("panic", r))
		}
	}()

	cacheUsers := make([]dao.User, len(users))

	for i, user := range users {
		cacheUsers[i] = dao.FromEntity(user)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	err := u.redis.SetCache(ctx, KeyAllUsers, cacheUsers, time.Hour)
	if err != nil {
		u.logger.Warn("redis set error", zap.Error(err))
	}
}

func (u *UserRepo) Store(ctx context.Context, user entity.User) error {
	err := u.repo.Store(ctx, user)
	if err != nil {
		return err
	}

	if err = u.redis.DeleteCache(ctx, KeyAllUsers); err != nil {
		u.logger.Error("redis del error", zap.String("key", KeyAllUsers), zap.Error(err))
	}

	return nil
}
