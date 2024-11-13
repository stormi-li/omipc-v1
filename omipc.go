package omipc

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi-v1"
)

func NewClient(opts *redis.Options) *Client {
	return &Client{
		redisClient:   redis.NewClient(opts),
		ctx:           context.Background(),
		configManager: omi.NewConfigManager(opts),
	}
}
