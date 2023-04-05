package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/stats", getStats).Methods("GET")
	r.HandleFunc("/ws", wsEndpoint)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	// Read the polling interval from the environment variable
	envPollingInterval := os.Getenv("POLLING_INTERVAL")
	if envPollingInterval == "" {
		log.Fatal("POLLING_INTERVAL environment variable not set")
	}

	// Parse the polling interval as a duration
	pollingInterval, err := time.ParseDuration(envPollingInterval)
	if err != nil {
		log.Fatalf("Failed to parse POLLING_INTERVAL: %v", err)
	}

	pollPlayerStats(pollingInterval)
}
