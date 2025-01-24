package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Leo7Deng/ChatApp/middleware"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/postgres"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(circles)
}

func CreateCirclesHandler(w http.ResponseWriter, r *http.Request) {
	var circle models.CircleData
	err := json.NewDecoder(r.Body).Decode(&circle)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("HTTP 400 bad request\n")
		return
	}
	userID := r.Context().Value(middleware.UserIDKey).(string)
	fmt.Println("Got userID from create circles: " + userID)
	err = postgres.CreateCircle(userID, circle.Name)
	if err != nil {
		fmt.Printf("Failed to create circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to create circle\n")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Circle created\n")
}