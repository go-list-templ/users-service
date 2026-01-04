package redis

import (
	"context"

	"github.com/go-list-templ/grpc/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func New(cfg *config.Redis) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

func (r *Redis) Invalidate(key string) error {
	return r.Del(context.Background(), key).Err()
}
