package auth

import (
	"fmt"
	"github.com/Leo7Deng/ChatApp/postgres"
	"github.com/Leo7Deng/ChatApp/redis"
	"github.com/google/uuid"
)

func CreateRefreshToken(accountEmail string) string {
	refreshToken := uuid.New()
    fmt.Println(accountEmail + " : " + refreshToken.String())
	postgres.InsertRefreshToken(accountEmail, refreshToken)
	redis.InsertRefreshToken(accountEmail, refreshToken.String())
	return refreshToken.String()
}