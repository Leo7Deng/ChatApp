package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/postgres"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var account user.FormData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		// return HTTP 400 bad request
		fmt.Printf("HTTP 400 bad request")
	} else {
		fmt.Printf("Registered account with email: %s\n", account.Email)
	}
	postgres.CreateAccount(account)
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var account user.FormData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		// return HTTP 400 bad request
		fmt.Printf("HTTP 400 bad request")
	} else {
		fmt.Printf("Logged in account with email: %s\n", account.Email)
	}
	// isLoggedIn := postgres.Login(account)
	isLoggedIn := true
	if isLoggedIn {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
	fmt.Println("LoginHandler")
}