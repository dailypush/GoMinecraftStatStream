package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func getStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := fetchPlayerStats(ctx) // Assuming fetchPlayerStats now accepts a context
	if err != nil {
		log.Printf("Error fetching player stats: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding stats to JSON: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func getAllStatTypesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	allStatTypes, err := getAllStatTypes(ctx)
	if err != nil {
		log.Printf("Error getting all stat types: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allStatTypes); err != nil {
		log.Printf("Error encoding all stat types to JSON: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func getSummarizedStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the statTypes parameter from the request
	ctx := r.Context()
	statTypesParam := r.URL.Query().Get("statType")
	if statTypesParam == "" {
		http.Error(w, "statType is required", http.StatusBadRequest)
		return
	}
	// Add further processing here...
}
