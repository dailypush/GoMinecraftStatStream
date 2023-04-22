package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type QueryParams struct {
	PlayerNames []string
	StatType    string
	GroupBy     string
	SortOrder   string
	Top         int
	Category    string
}

func parseQueryParams(r *http.Request) (QueryParams, error) {
	playerNamesStr := r.URL.Query().Get("playernames")
	log.Printf("Processing request for player names: %s", playerNamesStr)

	var playerNames []string
	if playerNamesStr != "" {
		playerNames = strings.Split(playerNamesStr, ",")
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		return QueryParams{}, errors.New("invalid sort order: must be either 'asc' or 'desc'")
	}

	groupBy := r.URL.Query().Get("groupby")
	if groupBy != "" && groupBy != "stattype" {
		return QueryParams{}, errors.New("invalid groupby option")
	}

	top := r.URL.Query().Get("top")
	topInt := 0
	if top != "" {
		var err error
		topInt, err = strconv.Atoi(top)
		if err != nil {
			return QueryParams{}, errors.New("invalid top value")
		}
	}
	category := r.URL.Query().Get("category")

	statType := r.URL.Query().Get("stattype")
	statType = strings.ReplaceAll(statType, "-", ":")

	return QueryParams{
		PlayerNames: playerNames,
		StatType:    statType,
		GroupBy:     groupBy,
		SortOrder:   sortOrder,
		Top:         topInt,
		Category:    category,
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
	log.Printf("Query params: %+v\n", queryParams)

	fmt.Println("Fetching player stats...")
	var playerStats []PlayerStats

	var allPlayerStats []PlayerStats
	for _, playerName := range queryParams.PlayerNames {
		redisPattern := fmt.Sprintf("player_stats:%s:*", playerName)
		keys, err := rdb.Keys(ctx, redisPattern).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stats, err := getPlayerStatsFromKeys(ctx, keys, playerName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		allPlayerStats = append(allPlayerStats, stats...)
	}

	playerStats = allPlayerStats
	if len(queryParams.PlayerNames) == 0 {
		redisPattern := "player_stats:*"

		keys, err := rdb.Keys(ctx, redisPattern).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(keys) == 0 {
			fetchPlayerStatsFromJson()
			keys, err = rdb.Keys(ctx, redisPattern).Result()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		fmt.Printf("Found %d keys in Redis.\n", len(keys))
		log.Printf("Keys: %+v\n", keys)
		playerStats, err = getPlayerStatsFromKeys(ctx, keys, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if queryParams.GroupBy == "stattype" {
		playerStats = groupByStatType(playerStats)
	}

	if queryParams.SortOrder != "" {
		sortByValue(playerStats, queryParams.SortOrder)
	}
	if queryParams.Top > 0 && queryParams.Category != "" {
		playerStats = getTopCategoryItems(playerStats, queryParams.Category, queryParams.Top)
	}
	log.Printf("Returning %d player stats\n", len(playerStats))
	writeJSONResponse(w, playerStats)
}

func getTopCategoryItems(playerStats []PlayerStats, category string, top int) []PlayerStats {
	filteredStats := []PlayerStats{}
	lowerCategory := strings.ToLower(category)

	for _, stat := range playerStats {
		log.Printf("StatType: %s", stat.StatType)
		if strings.HasPrefix(strings.ToLower(stat.StatType), lowerCategory) {
			filteredStats = append(filteredStats, stat)
		}
	}

	log.Printf("Filtered stats by category '%s': %d stats", category, len(filteredStats))
	for _, filteredStat := range filteredStats {
		log.Printf("Filtered stat: %+v", filteredStat)
	}

	sortByValue(filteredStats, "desc")

	if top > 0 && top < len(filteredStats) {
		filteredStats = filteredStats[:top]
	}
	log.Printf("Returning top %d items in category '%s'\n", top, category)
	return filteredStats
}

func getPlayerStatsFromKeys(ctx context.Context, keys []string, playerName string) ([]PlayerStats, error) {
	var playerStats []PlayerStats

	// Add this log statement to print the number of keys found in Redis
	log.Printf("Found %d keys in Redis.", len(keys))
	for _, key := range keys {
		value, err := rdb.Get(ctx, key).Int64()
		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %q: %v", key, err)
		}

		parts := strings.Split(key, ":")
		if len(parts) < 3 {
			continue
		}
		playerName := parts[1]
		statType := strings.Join(parts[2:], ":")

		playerStat := PlayerStats{
			Player:   playerName,
			StatType: statType,
			Value:    int(value),
		}
		playerStats = append(playerStats, playerStat)
		//log.Printf("Retrieved player stat from Redis: Key=%s, Value=%d", key, value)
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
