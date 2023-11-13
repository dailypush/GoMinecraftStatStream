package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// GetCurrentPlayers handles the HTTP request to retrieve current players.
// It sends a JSON response with the list of players or an error message.
func GetCurrentPlayers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Retrieve the context from the HTTP request
	players, err := getAllPlayersFromRedis(ctx)
	if err != nil {
		log.Printf("Error retrieving players from Redis: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(players) == 0 {
		http.Error(w, "No players found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(players); err != nil {
		log.Printf("Error encoding players to JSON: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// getAllPlayersFromRedis retrieves a list of current players from Redis.
// It accepts a context and returns an error if the retrieval or processing fails.
func getAllPlayersFromRedis(ctx context.Context) ([]string, error) {
	keys, err := rdb.Keys(ctx, "player_stats:*").Result() // Use the provided context
	if err != nil {
		return nil, fmt.Errorf("failed to get player keys from Redis: %v", err)
	}

	playerMap := make(map[string]struct{})
	for _, key := range keys {
		player := strings.Split(key, ":")[1]
		playerMap[player] = struct{}{}
	}

	players := make([]string, 0, len(playerMap))
	for player := range playerMap {
		players = append(players, player)
	}

	return players, nil
}
