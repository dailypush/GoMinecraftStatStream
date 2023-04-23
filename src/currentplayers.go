package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetCurrentPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := getAllPlayersFromRedis()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func getAllPlayersFromRedis() ([]string, error) {
	keys, err := rdb.Keys(ctx, "player_stats:*").Result()
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
