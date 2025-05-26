package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis() (*redis.Client, error) {
	redisUrl := os.Getenv("REDIS_URL")
	redisPsw := os.Getenv("REDIS_PASSWORD")

	if redisUrl == "" {
		return nil, fmt.Errorf("REDIS_URL does'n setted")
	}
	if redisPsw == "" {
		log.Println("password for redis is empty")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPsw,
		DB:       0,
	})

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	return rdb, nil
}
