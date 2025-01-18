package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/Leo7Deng/ChatApp/redis"
	"github.com/Leo7Deng/ChatApp/middleware"
	"github.com/Leo7Deng/ChatApp/auth"
	"github.com/Leo7Deng/ChatApp/dashboard"
	"github.com/Leo7Deng/ChatApp/postgres"
	"github.com/joho/godotenv"
)

type api struct {
	addr string
}

func main() {
	var err error
	err = godotenv.Load()
    if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
    }

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
	mux.HandleFunc("/api/register", middleware.AddCorsHeaders(auth.RegisterHandler))
	mux.HandleFunc("/api/login", middleware.AddCorsHeaders(auth.LoginHandler))
	mux.Handle("/api/dashboard",
    middleware.AddCorsHeaders(
        middleware.AuthMiddleware(
            http.HandlerFunc(dashboard.DashboardHandler),
        ),
    ),
)
	log.Fatal(srv.ListenAndServe())
}
