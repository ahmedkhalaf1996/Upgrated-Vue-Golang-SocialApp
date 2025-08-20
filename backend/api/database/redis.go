package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// InitRedis initializes Redis connection
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_URL"),      // e.g., "localhost:6379"
		Password:     os.Getenv("REDIS_PASSWORD"), // no password if empty
		DB:           0,                           // use default DB
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Connected to Redis successfully")
}

// CloseRedis closes Redis connection
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
