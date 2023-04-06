package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get the player's name from the query parameter "playername"
	playerName := r.URL.Query().Get("playername")
	if playerName == "" {
		http.Error(w, "Player name is required", http.StatusBadRequest)
		return
	}

	// Get the stat type from the query parameter "stattype" (optional)
	statType := r.URL.Query().Get("stattype")

	// Define a Redis pattern to match the stats for the given player
	redisPattern := "player_stats:" + playerName
	if statType != "" {
		redisPattern += ":" + statType
	} else {
		redisPattern += ":*"
	}

	// Get all keys matching the pattern
	keys, err := rdb.Keys(ctx, redisPattern).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve stats for each key and append to the result slice
	var playerStats []PlayerStats
	for _, key := range keys {
		value, err := rdb.Get(ctx, key).Int64()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		statType := strings.TrimPrefix(key, "player_stats:"+playerName+":") // Extract the stat type from the key
		playerStat := PlayerStats{
			Player:   playerName,
			StatType: statType,
			Value:    int(value),
		}
		playerStats = append(playerStats, playerStat)
	}

	// Convert the result slice to JSON and respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playerStats)
}
