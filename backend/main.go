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
	pool := postgres.GetPool()
	defer pool.Close()

	cassandraSession := cassandra.CassandraSession()
	defer cassandraSession.Close()

	hub := websockets.NewHub()
	
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	kafka.InitKafka(ctx, &wg, hub, cassandraSession, pool)

	go hub.Run() 
	fmt.Println("Websocket server started", hub)

	mux.HandleFunc("/api/register", middleware.AddCorsHeaders(auth.RegisterHandler))
	mux.HandleFunc("/api/login", middleware.AddCorsHeaders(auth.LoginHandler))
	mux.HandleFunc("/api/circles/search", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				circles.SearchTextHandler(w, r)
			}),
		),
	))
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
	mux.Handle("/api/circles/edit", middleware.AddCorsHeaders(
		middleware.AuthMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				circles.EditUsersHandler(w, r)
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
		// log.Fatal(srv.ListenAndServeTLS("full-cert.crt", "private-key.key"))
		log.Fatal(srv.ListenAndServe())
	}()

	wg.Wait()
	cancel()
}
