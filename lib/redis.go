package lib

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func SetCache(key string, value string, ttlMinutes int) bool {
	var redis_host = os.Getenv("REDIS_HOST")
	var redis_password = os.Getenv("REDIS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host,
		Password: redis_password,
		DB:       0,
	})
	if rdb == nil {
		panic("Failed to create Redis client")
	}

	println("Created Client in Set")

	err := rdb.Set(context.Background(), key, value, time.Minute*time.Duration(ttlMinutes)).Err()
	if err != nil {
		panic(err)
	}

	return err != nil
}

func GetCache(key string) string {
	var redis_host = os.Getenv("REDIS_HOST")
	var redis_password = os.Getenv("REDIS_PASSWORD")
	println("Entered Get")

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
		return "nil"
	}
	println("Called Get")

	return val
}
