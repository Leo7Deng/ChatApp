package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FormData struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var account FormData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		// return HTTP 400 bad request
		fmt.Printf("HTTP 400 bad request")
	} else {
		fmt.Printf("First name is %s\n", account.FirstName)
	}
}