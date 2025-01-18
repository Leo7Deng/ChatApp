package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start of dashboard handler")
	var data string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println("Failed to decode dashboard data")
	} else {
		fmt.Println("Decoded dashboard data: " + data)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Logged in\n")
}