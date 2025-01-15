package auth

import (
	"encoding/json"
	"fmt"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/postgres"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var account user.RegisterData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
	} else {
		fmt.Printf("Registered account with email: %s\n", account.Email)
	}
	var userID int
	userID, err = postgres.CreateAccount(account)
	if err != nil {
		fmt.Printf("Account creation failed\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Account creation failed\n")
	} else {
		fmt.Printf("Account created successfully\n")

		refreshToken, err := CreateRefreshToken(userID)
		if err != nil {
			fmt.Printf("Failed to create refresh token\n")
		}
		setRefreshTokenCookie(w, refreshToken)

		accessToken, err := CreateAccessToken(userID)
		if err != nil {
			fmt.Printf("Failed to create access token\n")
		}
		setAccessTokenCookie(w, accessToken)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Account created\n")
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var account user.LoginData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
	} else {
		fmt.Printf("Logged in account with email: %s\n", account.Email)
	}
	var userID int
	userID, err = postgres.FindAccount(account)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("Unauthorized login\n")
	} else {
		fmt.Printf("Logged in successfully\n")

		refreshToken, err := CreateRefreshToken(userID)
		if err != nil {
			fmt.Printf("Failed to create refresh token\n")
		}
		setRefreshTokenCookie(w, refreshToken)

		accessToken, err := CreateAccessToken(userID)
		if err != nil {
			fmt.Printf("Failed to create access token\n")
		}
		setAccessTokenCookie(w, accessToken)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Logged in\n")
	}
}

func setRefreshTokenCookie(w http.ResponseWriter, token string) {
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh-token",
        Value:    token,
        Path:     "/",
        Secure:   true,
        HttpOnly: true,
        MaxAge:   30 * 24 * 60 * 60, // 30 days
    })
}

func setAccessTokenCookie(w http.ResponseWriter, token string) {
    http.SetCookie(w, &http.Cookie{
        Name:     "access-token",
        Value:    token,
        Path:     "/",
        Secure:   true,
        HttpOnly: true,
        MaxAge:   15 * 60, // 15 minutes
    })
}