package auth

import (
    "fmt"
	"github.com/Leo7Deng/ChatApp/postgres"
    "github.com/google/uuid"
)

func CreateRefreshToken(accountEmail string) {
	refreshToken := uuid.New()
    fmt.Println(refreshToken.String())
	postgres.InsertRefreshToken(accountEmail, refreshToken)
}