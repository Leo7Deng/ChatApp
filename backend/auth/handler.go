package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/Leo7Deng/ChatApp/models"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var account user.FormData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		// return HTTP 400 bad request
		fmt.Printf("HTTP 400 bad request")
	} else {
		fmt.Printf("First name is %s\n", account.FirstName)
	}
}