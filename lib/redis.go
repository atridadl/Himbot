package lib

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redis_host = os.Getenv("REDIS_HOST")
var redis_password = os.Getenv("REDIS_PASSWORD")

func SetCache(key string, value string, ttlMinutes int) bool {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host,
		Password: redis_password,
		DB:       0,
	})

	err := rdb.Set(ctx, key, value, time.Minute*time.Duration(ttlMinutes)).Err()

	return err != nil
}

func GetCache(key string) string {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host,
		Password: redis_password,
		DB:       0,
	})

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "nil"
	}

	return val
}
