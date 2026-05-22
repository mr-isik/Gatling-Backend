package infra

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg RedisConfig) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error parsing redis url: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("error pinging redis: %w", err)
	}

	return client, nil
}
