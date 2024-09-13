package utils

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ConnectRedis(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	return rdb
}
