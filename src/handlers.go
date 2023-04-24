package main

import (
	"encoding/json"
	"net/http"
	"strings"
	// "github.com/gorilla/websocket"
)

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

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
	// Parse the statTypes parameter from the request
	statTypesParam := r.URL.Query().Get("statType")
	if statTypesParam == "" {
		http.Error(w, "statType is required", http.StatusBadRequest)
		return
	}

	// Split the statTypes parameter into a slice of stat types
	statTypes := strings.Split(statTypesParam, ",")

	// Get the summarized stats for each stat type and store them in a map
	aggregatedStats, individualStats, err := getSummarizedStats(statTypes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response JSON
	response := map[string]interface{}{
		"aggregatedStats": aggregatedStats,
		"individualStats": individualStats,
	}

	// Send the aggregated stats as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func wsEndpoint(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer ws.Close()

// 	for {
// 		_, _, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Println("read:", err)
// 			break
// 		}

// 		stats := fetchPlayerStats()
// 		err = ws.WriteJSON(stats)
// 		if err != nil {
// 			log.Println("write:", err)
// 			break
// 		}
// 	}
// }
