package auth

import (
	"fmt"
	"os"
	"strconv"
	"github.com/Leo7Deng/ChatApp/postgres"
	"github.com/Leo7Deng/ChatApp/redis"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func CreateRefreshToken(userID int) (string, error) {
	refreshToken := uuid.New()
	fmt.Println("User " + strconv.Itoa(userID) + "'s Refresh Token: " + refreshToken.String())
	postgres.InsertRefreshToken(userID, refreshToken)
	redis.InsertRefreshToken(userID, refreshToken.String())
	return refreshToken.String(), nil
}

func CreateAccessToken(userID int) (string, error) {
	var (
		key []byte
		t   *jwt.Token
		s   string
		err error
	  )
	  
	  key = []byte(os.Getenv("TOKEN_SECRET_KEY"))
	  t = jwt.New(jwt.SigningMethodHS256) 
	  s, err = t.SignedString(key)

	  if err != nil {
		return "", err
	  } else {
		fmt.Println("User " + strconv.Itoa(userID) + "'s Access Token: " + s)
		return s, nil
	  }
}