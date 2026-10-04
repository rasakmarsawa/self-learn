package connection

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return redisClient
}