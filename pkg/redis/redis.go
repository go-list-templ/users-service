package redis

import (
	"context"
	"encoding/json"
	"time"

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

func (r *Redis) DeleteCache(key string) error {
	return r.Del(context.Background(), key).Err()
}

func (r *Redis) GetCache(key string, pointer any) error {
	data, err := r.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, pointer)
}

func (r *Redis) SetCache(key string, data any, ttl time.Duration) error {
	data, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.Set(context.Background(), key, data, ttl).Err()
}
