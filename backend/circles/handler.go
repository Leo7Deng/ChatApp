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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	username, err := postgres.GetUsername(userID)
	if err != nil {
		fmt.Printf("Failed to get username\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to get username\n")
		return
	}
	w.WriteHeader(http.StatusOK)
	response := struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
	}{
		UserID:   userID,
		Username: username,
	}
	json.NewEncoder(w).Encode(response)
}

func GetInviteUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	type Circle struct {
		ID string `json:"circle_id"`
	}
	var circle Circle
	err := json.NewDecoder(r.Body).Decode(&circle)
	fmt.Println(json.NewDecoder(r.Body))
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("HTTP 400 bad request")
		return
	}
	fmt.Println("Circle ID INVITE: " + circle.ID + " USER ID: " + userID)
	users, err := postgres.GetInviteUsersInCircle(userID, circle.ID)
	if err != nil {
		fmt.Printf("Failed to get users in circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to get users in circle")
		return
	}
	fmt.Printf("Users in circle: %v\n", users)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func EditUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		GetEditUsersHandler(w, r)
	case "PUT":
		UpdateEditUsersHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetEditUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	type Circle struct {
		ID string `json:"circle_id"`
	}
	var circle Circle
	err := json.NewDecoder(r.Body).Decode(&circle)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("HTTP 400 bad request")
		return
	}
	users, err := postgres.GetExistingUsersInCircle(userID, circle.ID)
	if err != nil {
		fmt.Printf("Failed to get users in circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to get users in circle")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func UpdateEditUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	type EditData struct {
		ID     string `json:"circle_id"`
		UserID []struct {
			ID   string `json:"id"`
			Role string `json:"role"`
		} `json:"users"`
	}
	var circle EditData
	err := json.NewDecoder(r.Body).Decode(&circle)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("HTTP 400 bad request")
		return
	}
	for _, user := range circle.UserID {
		fmt.Println("Circle ID EDIT: " + circle.ID + " USER ID: " + userID + " EDIT USER ID: " + user.ID + " ROLE: " + user.Role)
		err = postgres.EditRoleInCircle(circle.ID, user.ID, user.Role)
	}
	if err != nil {
		fmt.Printf("Failed to edit users in circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to edit users in circle")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Users edited in circle")
}

func CreateCircleHandler(w http.ResponseWriter, r *http.Request, hub *websockets.Hub) {
	var createCircleData models.CreateCircleData
	err := json.NewDecoder(r.Body).Decode(&createCircleData)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("HTTP 400 bad request")
		return
	}
	userID := r.Context().Value(middleware.UserIDKey).(string)
	fmt.Println("Got userID from create circles: " + userID)

	var circle models.Circle
	circle, err = postgres.CreateCircle(userID, createCircleData.Name)
	if err != nil {
		fmt.Printf("Failed to create circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to create circle")
		return
	}
	hub.AddUsersToCircle(circle.ID, []string{userID})

	websocketMessage := models.WebsocketMessage{
		Origin: "server",
		Type: "circle",
		Action: "create",
		Message: nil,
		Circle: &circle,
	}
	hub.SendWebsocketMessage(websocketMessage)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Circle created")
}

func AddUsersToCircleHandler(w http.ResponseWriter, r *http.Request, hub *websockets.Hub) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	type AddData struct {
		ID     string   `json:"circle_id"`
		UserID []string `json:"users"`
	}
	fmt.Println("User " + userID + " is adding " + fmt.Sprint(userID) + " to circle")
	var circle AddData
	err := json.NewDecoder(r.Body).Decode(&circle)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("HTTP 400 bad request")
		return
	}
	fmt.Println("Circle ID: " + circle.ID + " User IDs: " + fmt.Sprint(circle.UserID))

	err = postgres.AddUsersToCircle(circle.ID, circle.UserID)
	if err != nil {
		fmt.Printf("Failed to add users to circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to add users to circle")
		return
	}

	hub.AddUsersToCircle(circle.ID, circle.UserID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Users added to circle")
}

func DeleteCircleHandler(w http.ResponseWriter, r *http.Request, hub *websockets.Hub) {
	circleID := strings.TrimPrefix(r.URL.Path, "/api/circles/delete/")
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
		if err.Error() == "permission error" {
			fmt.Printf("User is not admin of circle\n")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("user is not admin of circle")
			return
		}
		fmt.Printf("Failed to delete circle\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to delete circle\n")
		return
	}

	websocketMessage := models.WebsocketMessage{
		Origin: "server",
		Type: "circle",
		Action: "delete",
		Message: nil,
		Circle: &models.Circle{ID: circleID},
	}
	hub.SendWebsocketMessage(websocketMessage)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Circle deleted\n")
}

func SearchTextHandler(w http.ResponseWriter, r *http.Request) {
	type SearchData struct {
		CircleID string `json:"circle_id"`
		Content string `json:"content"`
	}
	var searchData SearchData
	err := json.NewDecoder(r.Body).Decode(&searchData)
	if err != nil {
		fmt.Printf("HTTP 400 bad request\n")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("HTTP 400 bad request")
		return
	}
	fmt.Println("Circle ID: " + searchData.CircleID + " Content: " + searchData.Content)
	results, err := postgres.SearchCircle(searchData.CircleID, searchData.Content)
	if err != nil {
		fmt.Printf("Failed to search circles\n")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to search circles\n")
		return
	}
	type Stats struct {
		TotalMatches int `json:"total_matches"`
		AverageLength float32 `json:"average_length"`
		MostFoundAuthor string `json:"most_found_author"`
	}
	var stats Stats
	stats.TotalMatches = len(results)
	if stats.TotalMatches > 0 {
		var totalLength int
		authorMap := make(map[string]int)
		for _, result := range results {
			totalLength += len(result.Content)
			authorMap[result.AuthorUsername]++
		}
		stats.AverageLength = float32(totalLength) / float32(stats.TotalMatches)
		var maxCount int
		var mostFoundAuthor string
		for author, count := range authorMap {
			if count > maxCount {
				maxCount = count
				mostFoundAuthor = author
			}
		}
		stats.MostFoundAuthor = mostFoundAuthor
	}

	fmt.Printf("Search stats: %v\n", stats)

	fmt.Printf("Search results: %v\n", results)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Results []models.SearchMessage `json:"results"`
		Stats   Stats                  `json:"stats"`
	}{
		Results: results,
		Stats:   stats,
	})
}