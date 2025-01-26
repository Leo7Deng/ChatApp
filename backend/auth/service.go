package auth

import (
	"fmt"
	"os"
	"time"
	"github.com/Leo7Deng/ChatApp/postgres"
	"github.com/Leo7Deng/ChatApp/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateRefreshToken(userID string) (string, error) {
	refreshToken := uuid.New()
	fmt.Println("User " + userID + "'s Refresh Token: " + refreshToken.String())
	postgres.InsertRefreshToken(userID, refreshToken)
	redis.InsertRefreshToken(userID, refreshToken.String())
	return refreshToken.String(), nil
}

func CreateAccessToken(userID string) (string, error) {
	secretKey := []byte(os.Getenv("TOKEN_SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	} else {
		fmt.Println("User " + userID + "'s Access Token: " + tokenString)
		return tokenString, nil
	}
}

func ValidateAccessToken(tokenString string) (string, error) {
	secretKey := []byte(os.Getenv("TOKEN_SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	} else {
		return "", err
	}
}