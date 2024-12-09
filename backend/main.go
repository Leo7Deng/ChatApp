package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/Leo7Deng/ChatApp/redis"
	"github.com/Leo7Deng/ChatApp/middleware"
	"github.com/Leo7Deng/ChatApp/auth"
	"github.com/Leo7Deng/ChatApp/postgres"
)

type api struct {
	addr string
}

func addCorsHeader(w http.ResponseWriter) {
    headers := w.Header()
    headers.Add("Access-Control-Allow-Origin", "*")
    headers.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	headers.Add("Access-Control-Allow-Headers", "Content-Type")
}



func main() {
	fmt.Println("Backend listening on port 8000")
	api := &api{addr: ":8000"}
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr: api.addr,
		Handler: mux,
	}
	redis.RedisClient()
	postgres.ConnectPSQL()
	defer postgres.ClosePSQL()
	// CORS for development
	mux.HandleFunc("/api/create_account", middleware.AddCorsHeaders(auth.Handler))
	log.Fatal(srv.ListenAndServe())
	
}
