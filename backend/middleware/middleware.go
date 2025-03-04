package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Leo7Deng/ChatApp/auth"
)

func AddCorsHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}
		if r.Header.Get("Origin") == "https://leo7deng.github.io" {
			w.Header().Set("Access-Control-Allow-Origin", "https://leo7deng.github.io")
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, GET, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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
		userID, err := auth.ValidateAccessToken(accessToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		fmt.Println("Authenticated user with ID: " + userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
