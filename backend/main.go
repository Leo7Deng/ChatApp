package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Leo7Deng/ChatApp/auth"
	"github.com/Leo7Deng/ChatApp/circles"
	"github.com/Leo7Deng/ChatApp/kafka"
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
	
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	kafka.InitKafka(ctx, &wg, hub)

	go hub.Run() 
	fmt.Println("Websocket server started", hub)

	// Kafka setup


	// client, ctx := kafka.KafkaClient()
	// defer client.Close()

	// kafka.ConsumeMessages(client, ctx)

	// go func() {
	// 	for {
	// 		kafka.ProduceMessage(client, ctx, "foo", "Hello from the server")
	// 		<-time.After(5 * time.Second)
	// 	}
	// }()

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
	mux.Handle("/api/circles/invite", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				circles.GetInviteUsersHandler(w, r)
			}),
		),
	))
	mux.Handle("/api/circles/delete/{id}", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				circles.DeleteCircleHandler(w, r, hub)
			}),
		),
	))
	mux.Handle("/api/circles/invite/add", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				circles.AddUsersToCircleHandler(w, r, hub)
			}),
		),
	))
	mux.HandleFunc("/ws", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			websockets.ServeWs(hub),
		),
	))
	
	
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	wg.Wait()
	cancel()
}
