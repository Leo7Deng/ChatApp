package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net"
	"sync"

	"github.com/Leo7Deng/ChatApp/auth"
	"github.com/Leo7Deng/ChatApp/cassandra"
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

func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

func main() {
	fmt.Println("IP Address: ", GetOutboundIP())
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

	cassandra.CassandraClient()

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
			websockets.ServeWs(hub),
	))

	mux.HandleFunc("/api/user", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				circles.GetUserHandler(w, r)
			}),
		),
	))

	mux.HandleFunc("/refresh", middleware.AddCorsHeaders(auth.RefreshAccessTokenHandler))
	
	go func() {
		log.Println("Starting HTTPS server on https://127.0.0.1:8000")
		log.Fatal(srv.ListenAndServeTLS("full-cert.crt", "private-key.key"))
	}()

	wg.Wait()
	cancel()
}
