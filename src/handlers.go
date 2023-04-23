package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getStats(w http.ResponseWriter, r *http.Request) {
	stats := fetchPlayerStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func getAllStatTypesHandler(w http.ResponseWriter, r *http.Request) {
	allStatTypes, err := getAllStatTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allStatTypes)
}

func getSummarizedStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the statType parameter from the request
	statType := r.URL.Query().Get("statType")
	if statType == "" {
		http.Error(w, "statType is required", http.StatusBadRequest)
		return
	}

	// Get the aggregated stats and individual stats for the specified statType
	summarizedStats, individualStats, err := getSummarizedStats(statType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the aggregated stats and individual stats as a JSON response
	response := map[string]interface{}{
		"summarizedStats": summarizedStats,
		"individualStats": individualStats,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		stats := fetchPlayerStats()
		err = ws.WriteJSON(stats)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
