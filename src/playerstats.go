package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type QueryParams struct {
	PlayerName string
	StatType   string
	GroupBy    string
	SortOrder  string
}

func parseQueryParams(r *http.Request) (QueryParams, error) {
	playerName := r.URL.Query().Get("playername")
	if playerName == "" {
		return QueryParams{}, errors.New("player name is required")
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		return QueryParams{}, errors.New("invalid sort order")
	}

	groupBy := r.URL.Query().Get("groupby")
	if groupBy != "" && groupBy != "stattype" {
		return QueryParams{}, errors.New("invalid groupby option")
	}

	return QueryParams{
		PlayerName: playerName,
		StatType:   r.URL.Query().Get("stattype"),
		GroupBy:    groupBy,
		SortOrder:  sortOrder,
	}, nil
}

func sortByValue(playerStats []PlayerStats, order string) {
	sort.Slice(playerStats, func(i, j int) bool {
		if order == "desc" {
			return playerStats[i].Value > playerStats[j].Value
		}
		return playerStats[i].Value < playerStats[j].Value
	})
}

func GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	queryParams, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	redisPattern := "player_stats:" + queryParams.PlayerName
	if queryParams.StatType != "" {
		redisPattern += ":" + queryParams.StatType
	} else {
		redisPattern += ":*"
	}

	keys, err := rdb.Keys(ctx, redisPattern).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var playerStats []PlayerStats
	for _, key := range keys {
		value, err := rdb.Get(ctx, key).Int64()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get value for key %q: %v", key, err), http.StatusInternalServerError)
			return
		}
		statType := strings.TrimPrefix(key, "player_stats:"+queryParams.PlayerName+":")
		playerStat := PlayerStats{
			Player:   queryParams.PlayerName,
			StatType: statType,
			Value:    int(value),
		}
		playerStats = append(playerStats, playerStat)
	}

	if queryParams.GroupBy == "stattype" {
		groupedStats := make(map[string][]PlayerStats)
		for _, stat := range playerStats {
			groupedStats[stat.StatType] = append(groupedStats[stat.StatType], stat)
		}
		playerStats = []PlayerStats{}
		for _, group := range groupedStats {
			playerStats = append(playerStats, group...)
		}
	}

	if queryParams.SortOrder != "" {
		sortByValue(playerStats, queryParams.SortOrder)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(playerStats)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode JSON: %v", err), http.StatusInternalServerError)
	}
}
