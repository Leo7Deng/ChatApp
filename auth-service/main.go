package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type FormData struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func addCorsHeader(w http.ResponseWriter) {
    headers := w.Header()
    headers.Add("Access-Control-Allow-Origin", "*")
    headers.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	headers.Add("Access-Control-Allow-Headers", "Content-Type")
}

func handler(w http.ResponseWriter, r *http.Request) {
    // for CORS development
    addCorsHeader(w)
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

	var account FormData
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		// return HTTP 400 bad request
		fmt.Printf("HTTP 400 bad request")
	} else {
		fmt.Printf("First name is %s\n", account.FirstName)
	}
}

func main() {
	fmt.Println("Hi")
	http.HandleFunc("/api/create_account", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
