package main

import (
	"fmt"
	"github.com/Leo7Deng/ChatApp/auth"
	"github.com/Leo7Deng/ChatApp/dashboard"
	"github.com/Leo7Deng/ChatApp/middleware"
	"github.com/Leo7Deng/ChatApp/postgres"
	"github.com/Leo7Deng/ChatApp/redis"
	"github.com/Leo7Deng/ChatApp/websockets"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

type api struct {
	addr string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	fmt.Println("Backend listening on port 8000")
	api := &api{addr: ":8000"}
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}
	redis.RedisClient()
	defer redis.CloseRedis()

	postgres.ConnectPSQL()
	defer postgres.ClosePSQL()

	hub := websockets.NewHub()
	go hub.Run() 
	fmt.Println("Websocket server started", hub)

	mux.HandleFunc("/api/register", middleware.AddCorsHeaders(auth.RegisterHandler))
	mux.HandleFunc("/api/login", middleware.AddCorsHeaders(auth.LoginHandler))
	mux.Handle("/api/dashboard",
		middleware.AddCorsHeaders(
			middleware.AuthMiddleware(
				http.HandlerFunc(dashboard.DashboardHandler),
			),
		),
	)
	mux.Handle("/api/create-circle", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			http.HandlerFunc(dashboard.CreateCirclesHandler),
		),
	))
	mux.HandleFunc("/ws", middleware.AddCorsHeaders(websockets.ServeWs(hub)))
	log.Fatal(srv.ListenAndServe())
}
