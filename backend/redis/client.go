package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func RedisClient() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDIS_PASSWORD"), 
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
	ctx := context.Background()

	err = client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	_, err = client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected")
}

func InsertRefreshToken(userEmail string, refreshToken string) bool {
	ctx := context.Background()
	err := client.Set(ctx, refreshToken, userEmail, 30 * 24 * time.Hour).Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert into database: %v\n", err)
		return false
	}
	fmt.Println("Successfully inserted refresh token")
	return true
}