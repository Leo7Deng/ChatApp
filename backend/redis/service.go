package redis

import (
	"context"
	"time"
	"fmt"
	"os"
)

func InsertRefreshToken(userID string, refreshToken string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	err := client.Set(ctx, key, userID, 1 * 24 * time.Hour).Err() // 1 day
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert into Redis: %v\n", err)
		return false
	}
	fmt.Println("Successfully inserted refresh token")
	return true
}

func FindRefreshToken(refreshToken string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	userID, err := client.Get(ctx, key).Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find refresh token: %v\n", err)
		return "", err
	}
	fmt.Println("Successfully found refresh token")
	return userID, nil
}
