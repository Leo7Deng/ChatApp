package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	
	"github.com/golang-jwt/jwt/v5"
)

func AddCorsHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://leo7deng.github.io")
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
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
		accessToken := r.Header.Get("Authorization")[7:]
		
		// If access token not found, return unauthorized
		if accessToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
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
