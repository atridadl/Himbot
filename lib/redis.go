package lib

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var redis_host = os.Getenv("REDIS_HOST")
var redis_password = os.Getenv("REDIS_PASSWORD")

func SetCache(key string, value string, ttlMinutes int) bool {
	println("Setting the Cache")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host,
		Password: redis_password,
		DB:       0,
	})
	if rdb == nil {
		panic("Failed to create Redis client")
	}

	err := rdb.Set(context.Background(), key, value, time.Minute*time.Duration(ttlMinutes)).Err()
	if err != nil {
		panic(err)
	}

	return err != nil
}

func GetCache(key string) string {
	println("Fetching From Cache")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host,
		Password: redis_password,
		DB:       0,
	})
	if rdb == nil {
		panic("Failed to create Redis client")
	}

	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		println("Cache Miss")
		return "nil"
	}

	println("Cache Hit")

	return val
}
