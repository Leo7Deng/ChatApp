package auth

import (
	"fmt"
	"strconv"
	"github.com/Leo7Deng/ChatApp/postgres"
	"github.com/Leo7Deng/ChatApp/redis"
	"github.com/google/uuid"
	// "github.com/golang-jwt/jwt/v5"
)

func CreateRefreshToken(userID int) string {
	refreshToken := uuid.New()
	fmt.Println("User " + strconv.Itoa(userID) + "'s Refresh Token: " + refreshToken.String())
	postgres.InsertRefreshToken(userID, refreshToken)
	redis.InsertRefreshToken(userID, refreshToken.String())
	return refreshToken.String()
}

// func CreateAccessToken(accountEmail string) string {
// 	var (
// 		key []byte
// 		t   *jwt.Token
// 		s   string
// 	  )
	  
// 	  key = 
// 	  /* Load key from somewhere, for example an environment variable */
// 	//   add user id as jwt claim
// 	  t = jwt.New(jwt.SigningMethodHS256) 
// 	  s = t.SignedString(key)
// }