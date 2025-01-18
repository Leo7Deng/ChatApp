package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/Leo7Deng/ChatApp/middleware"

)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start of dashboard handler")
	// userID := r.Context().Value(middleware.UserIDKey).(string)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Logged in\n")
}