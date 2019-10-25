package redis

import (
	redis "github.com/go-redis/redis/v7"
)

type client struct {
	redisClient *redis.Client
}
