package repo

import (
	"context"
	"errors"
	"time"

	"github.com/go-list-templ/grpc/internal/adapter/cache/redis"
	"github.com/go-list-templ/grpc/internal/adapter/cache/redis/repo/dao"
	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/port"
	"github.com/go-list-templ/grpc/pkg/paginate"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var ErrTypedSingleflight = errors.New("invalid type from singleflight")

const (
	TTLAllUsers = 10 * time.Minute

	TagAllUsers = "allUsers"
)

type UserRepo struct {
	repo   port.UserRepo
	redis  *redis.Redis
	logger *zap.Logger
	sf     singleflight.Group
}

func NewUserRepo(repo port.UserRepo, redis *redis.Redis, logger *zap.Logger) *UserRepo {
	return &UserRepo{repo: repo, redis: redis, logger: logger, sf: singleflight.Group{}}
}

func (u *UserRepo) All(ctx context.Context, paginate paginate.Paginate) ([]entity.User, error) {
	var cachedUsers []dao.User
	cacheKey := paginate.Cursor()

	err := u.redis.GetCache(ctx, cacheKey, &cachedUsers)
	if err == nil && len(cachedUsers) > 0 {
		users := make([]entity.User, len(cachedUsers))

		for i, user := range cachedUsers {
			users[i] = user.ToEntity()
		}

		return users, nil
	}

	v, err, _ := u.sf.Do(cacheKey, func() (interface{}, error) {
		users, err := u.repo.All(ctx, paginate)
		if err != nil {
			return nil, err
		}

		cacheUsers := make([]dao.User, len(users))

		for i, user := range users {
			cacheUsers[i] = dao.FromEntity(user)
		}

		if err = u.redis.SetByTags(ctx, cacheKey, cacheUsers, TTLAllUsers, TagAllUsers); err != nil {
			u.logger.Warn("redis set failed", zap.Error(err))
		}

		return users, nil
	})
	if err != nil {
		return nil, err
	}

	users, ok := v.([]entity.User)
	if !ok {
		u.logger.Error("singleflight typed", zap.Error(err), zap.Any("value", v))

		return nil, ErrTypedSingleflight
	}

	return users, nil
}

func (u *UserRepo) Store(ctx context.Context, user entity.User) error {
	err := u.repo.Store(ctx, user)
	if err != nil {
		return err
	}

	if err = u.redis.InvalidateTags(ctx, TagAllUsers); err != nil {
		u.logger.Warn("redis invalidate error", zap.Error(err))
	}

	return nil
}
