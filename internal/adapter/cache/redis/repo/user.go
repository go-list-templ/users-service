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

const DefaultContextTimeout = 5 * time.Second

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

	// todo check what working
	v, err, _ := u.sf.Do(cacheKey, func() (interface{}, error) {
		users, err := u.repo.All(ctx, paginate)
		if err != nil {
			return nil, err
		}

		go u.cacheAllUsers(cacheKey, users)

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

func (u *UserRepo) cacheAllUsers(cacheKey string, users []entity.User) {
	cacheUsers := make([]dao.User, len(users))

	for i, user := range users {
		cacheUsers[i] = dao.FromEntity(user)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	if err := u.redis.SetCache(ctx, cacheKey, cacheUsers, time.Hour); err != nil {
		u.logger.Warn("redis set error", zap.Error(err))
	}
}

func (u *UserRepo) Store(ctx context.Context, user entity.User) error {
	err := u.repo.Store(ctx, user)
	if err != nil {
		return err
	}

	//if err = u.redis.DeleteCache(ctx, KeyAllUsers); err != nil {
	//	u.logger.Error("redis del error", zap.String("key", KeyAllUsers), zap.Error(err))
	//}

	return nil
}
