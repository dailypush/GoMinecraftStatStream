package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type QueryParams struct {
	PlayerName  string
	PlayerNames string
	StatType    string
	GroupBy     string
	SortOrder   string
	Top         int
	Category    string
}

func parseQueryParams(r *http.Request) (QueryParams, error) {
	playerName := r.URL.Query().Get("playername")
	playerNames := r.URL.Query().Get("playernames")

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		return QueryParams{}, errors.New("invalid sort order: must be either 'asc' or 'desc")
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
		PlayerName:  playerName,
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
	fmt.Println("Fetching player stats...")
	var playerStats []PlayerStats

	if queryParams.PlayerNames != "" {
		playerNames := strings.Split(queryParams.PlayerNames, ",")
		var allPlayerStats []PlayerStats
		for _, playerName := range playerNames {
			redisPattern := fmt.Sprintf("player_stats:%s", playerName)
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
	} else {
		redisPattern := "player_stats:*"
		if queryParams.PlayerName != "" {
			redisPattern = fmt.Sprintf("player_stats:%s", queryParams.PlayerName)
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
		fmt.Printf("Found %d keys in Redis.\n", len(keys))
		playerStats, err = getPlayerStatsFromKeys(ctx, keys, queryParams.PlayerName)
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

	writeJSONResponse(w, playerStats)
}

func getTopCategoryItems(playerStats []PlayerStats, category string, top int) []PlayerStats {
	var categoryItems []PlayerStats
	categoryPrefix := fmt.Sprintf("minecraft:%s:", category)

	for _, stat := range playerStats {
		if strings.HasPrefix(stat.StatType, categoryPrefix) {
			categoryItems = append(categoryItems, stat)
		}
	}

	sort.Slice(categoryItems, func(i, j int) bool {
		return categoryItems[i].Value > categoryItems[j].Value
	})

	if top > len(categoryItems) {
		top = len(categoryItems)
	}

	return categoryItems[:top]
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
