package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Start(port string) {
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: s.routes(),
	}

	go func() {
		log.Printf("Server is listening on port: %s", port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/getstats", getStats)
	mux.HandleFunc("/ws", wsEndpoint)
	mux.HandleFunc("/playerstats", GetPlayerStats)
	return mux
}

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Set the default port number here
	}

	server := &Server{}
	server.Start(port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")
	server.Stop()
	log.Println("Server stopped")
}
