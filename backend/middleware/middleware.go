package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Leo7Deng/ChatApp/auth"
	"github.com/Leo7Deng/ChatApp/redis"
	"github.com/golang-jwt/jwt/v5"
)

func AddCorsHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		var accessToken string
		var refreshToken string
		for _, c := range cookies {
			if c.Name == "access-token" {
				fmt.Println("Found access-token: " + c.Value)
				accessToken = c.Value
				continue
			}
			if c.Name == "refresh-token" {
				fmt.Println("Found refresh-token: " + c.Value)
				refreshToken = c.Value
				continue
			}
		}

		// If access token not found, try to refresh
		if accessToken == "" {
			fmt.Println("Access token not found, will try to refresh")
			userID, err := refreshAccessToken(w, refreshToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode("refresh token not found")
				return
			}
			fmt.Println("Refreshed access token")
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Validate access token
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("TOKEN_SECRET_KEY")), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			fmt.Println("2")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("3")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			fmt.Println("4")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		fmt.Println("Authenticated user with ID: " + userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func refreshAccessToken(w http.ResponseWriter, refreshToken string) (string, error) {
	if refreshToken == "" {
		fmt.Println("Refresh token not found, user will need to login again")
		return "", fmt.Errorf("refresh token not found")
	}
	userID, err := redis.FindRefreshToken(refreshToken)
	if err != nil {
		fmt.Println("Failed to find refresh token")
		return "", fmt.Errorf("refresh token not found")
	}
	accessToken, err := auth.CreateAccessToken(userID)
	if err != nil {
		fmt.Println("Failed to create access token")
		return "", fmt.Errorf("failed to create access token")
	}
	auth.SetAccessTokenCookie(w, accessToken)
	return userID, nil
}
