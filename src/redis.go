package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
	// Read Redis connection settings from environment variables
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	// Set default values if not provided in environment variables
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	if redisDB == "" {
		redisDB = "0"
	}

	// Parse Redis database index as integer
	redisDBInt, err := strconv.Atoi(redisDB)
	if err != nil {
		log.Fatalf("Invalid Redis database index: %s", redisDB)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDBInt,
	})
}

func storePlayerStatInRedis(stat PlayerStats) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("player_stats:%s:%s", stat.Player, stat.StatType)
	err := rdb.Set(ctx, key, stat.Value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set stat in Redis: %v", err)
	}

	return nil
}
