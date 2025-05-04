package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

const MAX_CACHE = 20

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

func SetHash(client *redis.Client, key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetHash(client *redis.Client, key string) (string, error) {
	result, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return result, fmt.Errorf("hash can't be found by key = %s", key)
	}

	if len(result) == 0 {
		return result, fmt.Errorf("hash is empty by key = %s", key)
	}

	return result, nil
}

func CashTask(client *redis.Client, cacheSetKey string, cacheHashKey string, taskId string, taskData interface{}) error {
	timestamp := float64(time.Now().UnixNano())

	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		fmt.Printf("error while marshaling to json: %v", err)
		return err
	}

	client.ZAdd(context.Background(), cacheSetKey, &redis.Z{
		Score:  timestamp,
		Member: taskId,
	})

	client.HSet(context.Background(), cacheHashKey, taskId, taskJSON)

	client.ZRemRangeByRank(context.Background(), cacheSetKey, 0, -int64(MAX_CACHE)-1)

	existingTasks, err := client.ZRange(context.Background(), cacheSetKey, 0, -1).Result()
	if err != nil {
		fmt.Println("Ошибка при получении задач из Sorted Set:", err)
		return err
	}

	allTaskIDs, err := client.HKeys(context.Background(), cacheHashKey).Result()
	if err != nil {
		fmt.Println("Ошибка при получении ключей из Hash:", err)
		return err
	}

	contains := func(slice []string, item string) bool {
		for _, v := range slice {
			if v == item {
				return true
			}
		}
		return false
	}

	for _, taskID := range allTaskIDs {
		if !contains(existingTasks, taskID) {
			client.HDel(context.Background(), cacheHashKey, taskID)
		}
	}
	return nil
}
