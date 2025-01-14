package auth

import (
	"fmt"
	"github.com/Leo7Deng/ChatApp/postgres"
	"github.com/google/uuid"
)

func CreateRefreshToken(accountEmail string) string {
	refreshToken := uuid.New()
    fmt.Println(accountEmail + " : " + refreshToken.String())
	postgres.InsertRefreshToken(accountEmail, refreshToken)
	return refreshToken.String()
}