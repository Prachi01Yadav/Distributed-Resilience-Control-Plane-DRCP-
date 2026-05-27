package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{Client: rdb}
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
