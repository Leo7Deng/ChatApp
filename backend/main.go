package main

import (
	"fmt"
	"log"
	"net/http"
	// "time"

	"github.com/Leo7Deng/ChatApp/auth"
	"github.com/Leo7Deng/ChatApp/circles"
	"github.com/Leo7Deng/ChatApp/middleware"
	"github.com/Leo7Deng/ChatApp/postgres"
	"github.com/Leo7Deng/ChatApp/redis"
	"github.com/Leo7Deng/ChatApp/websockets"
	"github.com/joho/godotenv"
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

	// asyncronously broadcast message every 5 seconds
	// go func() {
	// 	for {
	// 		hub.Broadcast([]byte("Hello from the server"))
	// 		fmt.Println("Broadcasting message")
	// 		<-time.After(5 * time.Second)
	// 	}
	// }()

	mux.HandleFunc("/api/register", middleware.AddCorsHeaders(auth.RegisterHandler))
	mux.HandleFunc("/api/login", middleware.AddCorsHeaders(auth.LoginHandler))
	mux.Handle("/api/circles", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				circles.CircleHandler(w, r, hub)
			}),
		),
	))
	mux.Handle("/api/circles/{id}", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				circles.DeleteCircleHandler(w, r, hub)
			}),
		),
	))
	mux.HandleFunc("/ws", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			websockets.ServeWs(hub),
		),
	))
	log.Fatal(srv.ListenAndServe())
}
