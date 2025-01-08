package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/postgres"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var account user.RegisterData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		// return HTTP 400 bad request
		fmt.Printf("HTTP 400 bad request\n")
	} else {
		fmt.Printf("Registered account with email: %s\n", account.Email)
	}
	if postgres.CreateAccount(account) {
		fmt.Printf("Account created successfully\n")
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Printf("Account creation failed\n")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var account user.LoginData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		// return HTTP 400 bad request
		fmt.Printf("HTTP 400 bad request\n")
	} else {
		fmt.Printf("Logged in account with email: %s\n", account.Email)
	}
	isLoggedIn := postgres.FindAccount(account)
	if isLoggedIn {
		w.WriteHeader(http.StatusOK)
		fmt.Printf("Logged in successfully\n")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("Unauthorized login\n")
	}
}