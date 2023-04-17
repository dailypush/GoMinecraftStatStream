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

	statType := r.URL.Query().Get("stattype")
	statType = strings.ReplaceAll(statType, "-", ":")

	return QueryParams{
		PlayerName: playerName,
		StatType:   statType,
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

	redisPattern := fmt.Sprintf("player_stats:%s", queryParams.PlayerName)
	if queryParams.StatType != "" {
		redisPattern = fmt.Sprintf("%s:%s", redisPattern, queryParams.StatType)
	} else {
		redisPattern = fmt.Sprintf("%s:*", redisPattern)
	}

	keys, err := rdb.Keys(ctx, redisPattern).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the keys are empty, and if so, fetch player stats from JSON files
	if len(keys) == 0 {
		fetchPlayerStatsFromJson()
		// Re-query Redis for the keys after loading stats from JSON files
		keys, err = rdb.Keys(ctx, redisPattern).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	playerStats, err := getPlayerStatsFromKeys(ctx, keys, queryParams.PlayerName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if queryParams.GroupBy == "stattype" {
		playerStats = groupByStatType(playerStats)
	}

	if queryParams.SortOrder != "" {
		sortByValue(playerStats, queryParams.SortOrder)
	}

	writeJSONResponse(w, playerStats)
}

func getPlayerStatsFromKeys(ctx context.Context, keys []string, playerName string) ([]PlayerStats, error) {
	var playerStats []PlayerStats
	for _, key := range keys {
		value, err := rdb.Get(ctx, key).Int64()
		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %q: %v", key, err)
		}
		statType := strings.TrimPrefix(key, fmt.Sprintf("player_stats:%s:", playerName))
		playerStat := PlayerStats{
			Player:   playerName,
			StatType: statType,
			Value:    int(value),
		}
		playerStats = append(playerStats, playerStat)
	}
	return playerStats, nil
}

func groupByStatType(playerStats []PlayerStats) []PlayerStats {
	groupedStats := make(map[string][]PlayerStats)
	for _, stat := range playerStats {
		groupedStats[stat.StatType] = append(groupedStats[stat.StatType], stat)
	}
	playerStats = []PlayerStats{}
	for _, group := range groupedStats {
		playerStats = append(playerStats, group...)
	}
	return playerStats
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode JSON: %v", err), http.StatusInternalServerError)
	}
}
