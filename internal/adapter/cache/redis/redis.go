package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-list-templ/grpc/config"
	"github.com/redis/go-redis/v9"
)

const DefaultContextTimeout = 5 * time.Second

type Redis struct {
	*redis.Client
}

func New(cfg *config.Redis) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

func (r *Redis) DeleteCache(ctx context.Context, key string) error {
	return r.Del(ctx, key).Err()
}

func (r *Redis) GetCache(ctx context.Context, key string, pointer any) error {
	data, err := r.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, pointer)
}

func (r *Redis) SetCache(ctx context.Context, key string, data any, ttl time.Duration) error {
	data, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.Set(ctx, key, data, ttl).Err()
}
