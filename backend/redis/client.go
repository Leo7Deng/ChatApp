package redis

import (
	"context"
	"fmt"
	"os"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func RedisClient() {
	fmt.Println("Creating Redis client")
	client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: os.Getenv("REDIS_PASSWORD"), 
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
	ctx := context.Background()

	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}
	_, err = client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	_, err = client.Del(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected")
}

func CloseRedis() {
	err := client.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to close Redis connection: %v\n", err)
	}
	fmt.Println("Redis connection closed")
}

func DeleteTable() bool {
	ctx := context.Background()
	err := client.FlushDB(ctx).Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to delete table: %v\n", err)
		return false
	}
	fmt.Println("Successfully deleted table")
	return true
}