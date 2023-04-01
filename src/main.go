package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type PlayerStats struct {
	Player   string `json:"player"`
	StatType string `json:"stat_type"`
	Value    int    `json:"value"`
}

var stats []PlayerStats

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/stats", getStats).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func init() {
	// Connect to your Minecraft server using RCON
	conn, err := websocket.Dial("ws://your-minecraft-server-address:RCON-port", "", "http://localhost/")
	if err != nil {
		log.Fatal("Failed to connect to RCON:", err)
	}
	defer conn.Close()

	// Authenticate using RCON password
	err = conn.WriteMessage(websocket.TextMessage, []byte("/rcon your-rcon-password"))
	if err != nil {
		log.Fatal("Failed to authenticate:", err)
	}

	// Fetch player stats
	stats = fetchPlayerStats(conn)
}

func fetchPlayerStats(conn *websocket.Conn) []PlayerStats {
	// Replace this with your actual command to fetch player stats
	statsCommand := "/your-stats-command"

	err := conn.WriteMessage(websocket.TextMessage, []byte(statsCommand))
	if err != nil {
		log.Fatal("Failed to send stats command:", err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Failed to read stats message:", err)
	}

	var stats []PlayerStats
	err = json.Unmarshal(message, &stats)
	if err != nil {
		log.Fatal("Failed to unmarshal stats:", err)
	}

	return stats
}
