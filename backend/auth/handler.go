package auth

import (
	"encoding/json"
	"fmt"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/postgres"
	"net/http"
	"github.com/Leo7Deng/ChatApp/redis"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var account models.RegisterData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
	} else {
		fmt.Printf("Registered account with email: %s\n", account.Email)
	}
	var userID string
	userRepo := postgres.NewUserRepository(postgres.GetPool())
	userID, err = userRepo.CreateAccount(account)
	if err != nil {
		if err.Error() == "email" {
			fmt.Printf("Email already exists\n")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode("Email already exists")
		} else if err.Error() == "username" {
			fmt.Printf("Username already exists\n")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode("Username already exists")
		} else {
			fmt.Printf("Account creation failed\n")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("Account creation failed")
		}
	} else {
		fmt.Printf("Account created successfully\n")

		refreshToken, err := CreateRefreshToken(userID)
		if err != nil {
			fmt.Printf("Failed to create refresh token\n")
		}
		SetRefreshTokenCookie(w, refreshToken)

		accessToken, err := CreateAccessToken(userID)
		if err != nil {
			fmt.Printf("Failed to create access token\n")
		}
		AccessTokenResponse := models.AccessTokenResponse{AccessToken: accessToken}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AccessTokenResponse)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var account models.LoginData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
	} else {
		fmt.Printf("Logged in account with email: %s\n", account.Email)
	}
	var user *models.User
	userRepo := postgres.NewUserRepository(postgres.GetPool())
	user, err = userRepo.FindAccount(account.Email)
	userID := user.ID
	if err != nil || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("Unauthorized login\n")
	} else {
		fmt.Printf("Logged in successfully\n")
		refreshToken, err := CreateRefreshToken(userID)
		if err != nil {
			fmt.Printf("Failed to create refresh token\n")
		}
		SetRefreshTokenCookie(w, refreshToken)

		accessToken, err := CreateAccessToken(userID)
		if err != nil {
			fmt.Printf("Failed to create access token\n")
		}
		AccessTokenResponse := models.AccessTokenResponse{AccessToken: accessToken}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AccessTokenResponse)
	}
}

func SetRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    token,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   1 * 24 * 60 * 60, // 1 day
		SameSite: http.SameSiteNoneMode,
	})
}

func RefreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	var refreshToken string
	for _, c := range cookies {
		if c.Name == "refresh-token" {
			refreshToken = c.Value
			continue
		}
	}
	if refreshToken == "" {
		fmt.Println("Refresh token not found, user will need to login again")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("refresh token not found")
	}
	userID, err := redis.FindRefreshToken(refreshToken)
	if err != nil {
		fmt.Println("Failed to find refresh token")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("refresh token not found")
	}
	accessToken, err := CreateAccessToken(userID)
	if err != nil {
		fmt.Println("Failed to create access token")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to create access token")
	}
	AccessTokenResponse := models.AccessTokenResponse{AccessToken: accessToken}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AccessTokenResponse)
}