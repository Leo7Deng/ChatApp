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
	if postgres.CreateAccount(account) {
		fmt.Printf("Account created successfully\n")
		refreshToken := CreateRefreshToken(account.Email)
		cookie := http.Cookie{
			Name:     "refresh-token",
			Value:    refreshToken,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			MaxAge:   30 * 24 * 60 * 60, // 30 days
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Account created\n")
	} else {
		fmt.Printf("Account creation failed\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Account creation failed\n")
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
	isLoggedIn := postgres.FindAccount(account)
	if isLoggedIn {
		fmt.Printf("Logged in successfully\n")
		refreshToken := CreateRefreshToken(account.Email)
		cookie := http.Cookie{
			Name:     "refresh-token",
			Value:    refreshToken,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			MaxAge:   30 * 24 * 60 * 60, // 30 days
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Logged in\n")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("Unauthorized login\n")
	}
}
