package circles

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leo7Deng/ChatApp/middleware"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/postgres"
	"github.com/Leo7Deng/ChatApp/websockets"
)

func CircleHandler(w http.ResponseWriter, r *http.Request, hub *websockets.Hub) {
	switch r.Method {
	case "GET":
		GetCirclesHandler(w, r)
	case "POST":
		CreateCircleHandler(w, r, hub)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetCirclesHandler(w http.ResponseWriter, r *http.Request) {
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


func CreateCircleHandler(w http.ResponseWriter, r *http.Request, hub *websockets.Hub) {
	var circleData models.CircleData
	err := json.NewDecoder(r.Body).Decode(&circleData)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("HTTP 400 bad request\n")
		return
	}
	userID := r.Context().Value(middleware.UserIDKey).(string)
	fmt.Println("Got userID from create circles: " + userID)

	var circle models.Circle
	circle, err = postgres.CreateCircle(userID, circleData.Name)
	if err != nil {
		fmt.Printf("Failed to create circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to create circle\n")
		return
	}

	response := struct {
		Type string        `json:"type"`
		Data models.Circle `json:"data"`
	}{
		Type: "add-circle",
		Data: circle,
	}
	circleJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Failed to marshal circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to marshal circle\n")
		return
	}
	hub.SendToUser(userID, circleJSON)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Circle created\n")
}

func DeleteCircleHandler(w http.ResponseWriter, r *http.Request, hub *websockets.Hub) {
	circleID := strings.TrimPrefix(r.URL.Path, "/api/circles/")
	fmt.Println("Circle ID: " + circleID)

	if circleID == "" {
		fmt.Printf("HTTP 400 bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("HTTP 400 bad request\n")
		return
	}
	userID := r.Context().Value(middleware.UserIDKey).(string)
	fmt.Println("Got userID from delete circles: " + userID)

	err := postgres.DeleteCircle(userID, circleID)
	if err != nil {
		fmt.Printf("Failed to delete circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to delete circle\n")
		return
	}

	response := struct {
		Type string        `json:"type"`
		Data models.Circle `json:"data"`
	}{
		Type: "remove-circle",
		Data: models.Circle{ID: circleID},
	}

	circleJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Failed to marshal circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to marshal circle\n")
		return
	}
	hub.Broadcast(circleJSON)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Circle deleted\n")
}
