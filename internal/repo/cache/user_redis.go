package cache

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo"
	"github.com/go-list-templ/grpc/internal/repo/cache/dao"
	"github.com/go-list-templ/grpc/pkg/redis"
	"go.uber.org/zap"
)

const (
	KeyAllUsers = "users:all"
)

type UserRedis struct {
	repo   repo.UserRepo
	redis  *redis.Redis
	logger *zap.Logger
}

func NewUserRedis(repo repo.UserRepo, redis *redis.Redis, logger *zap.Logger) *UserRedis {
	return &UserRedis{repo: repo, redis: redis, logger: logger}
}

func (u *UserRedis) All(ctx context.Context) ([]entity.User, error) {
	var cachedUsers []dao.User

	err := u.redis.GetCache(KeyAllUsers, &cachedUsers)
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

	go func() {
		cacheUsers := make([]dao.User, len(users))

		for i, user := range users {
			cacheUsers[i] = dao.FromEntity(user)
		}

		err = u.redis.SetCache(KeyAllUsers, cacheUsers, time.Hour)
		if err != nil {
			u.logger.Error("redis set error", zap.Error(err))
		}
	}()

	return users, nil
}

func (u *UserRedis) Store(ctx context.Context, tx pgx.Tx, user entity.User) error {
	err := u.repo.Store(ctx, tx, user)
	if err != nil {
		return err
	}

	go func() {
		if err = u.redis.DeleteCache(KeyAllUsers); err != nil {
			u.logger.Error("redis del error", zap.Error(err))
		}
	}()

	return nil
}
