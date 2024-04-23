package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Context = context.Background()

func CreateRedisClient(dbno int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       dbno,
	})

	return client
}
