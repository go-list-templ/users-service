package cache

import (
	"context"
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

func (u *UserRedis) cacheAllUsers(users []entity.User) {
	defer func() {
		if r := recover(); r != nil {
			u.logger.Error("panic in cacheAllUsers", zap.Any("panic", r))
		}
	}()

	cacheUsers := make([]dao.User, len(users))

	for i, user := range users {
		cacheUsers[i] = dao.FromEntity(user)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.redis.SetCache(ctx, KeyAllUsers, cacheUsers, time.Hour)
	if err != nil {
		u.logger.Error("redis set error", zap.Error(err))
	}
}

func (u *UserRedis) Store(ctx context.Context, user entity.User) error {
	err := u.repo.Store(ctx, user)
	if err != nil {
		return err
	}

	go u.clearCache(KeyAllUsers)

	return nil
}

func (u *UserRedis) clearCache(keys ...string) {
	defer func() {
		if r := recover(); r != nil {
			u.logger.Error("panic in clearCache", zap.Any("panic", r))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, key := range keys {
		if err := u.redis.DeleteCache(ctx, key); err != nil {
			u.logger.Error("redis del error", zap.String("key", key), zap.Error(err))
		}
	}
}
