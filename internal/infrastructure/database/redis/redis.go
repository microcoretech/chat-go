package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, redisHost, password string, redisDB int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: password,
		DB:       redisDB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
