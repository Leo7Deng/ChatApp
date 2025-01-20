package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/Leo7Deng/ChatApp/middleware"
	"github.com/Leo7Deng/ChatApp/postgres"

)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start of dashboard handler")
	userID := r.Context().Value(middleware.UserIDKey).(string)
	fmt.Println("User ID: " + userID)
	circles, err := postgres.GetUserCircles(userID)
	if err != nil {
		fmt.Printf("Failed to get user circles\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to get user circles\n")
		return
	}
	fmt.Printf("User circles: %v\n", circles)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Logged in\n")
}