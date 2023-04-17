package main

import (
	"context"
	"encoding/json"
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
	// mux.HandleFunc("/ws", wsEndpoint) // Define wsEndpoint if needed
	mux.HandleFunc("/playerstats", GetPlayerStats) // Referencing GetPlayerStats from playerstats.go
	return mux
}

func getStats(w http.ResponseWriter, r *http.Request) {
	stats := fetchPlayerStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {

	server := &Server{}
	server.Start(ServerPort)

	// Start the pollPlayerStats function in a separate goroutine
	go pollPlayerStats(PollingInterval)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")
	server.Stop()
	log.Println("Server stopped")
}
