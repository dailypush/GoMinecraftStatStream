package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const playerStatsPattern = "player_stats:%s"

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

	if len(queryParams.PlayerNames) > 0 {
		for _, playerName := range queryParams.PlayerNames {
			redisPattern := fmt.Sprintf("player_stats:%s", playerName)
			fmt.Println("redis pattern: ", redisPattern)
			keys, err := rdb.Keys(ctx, redisPattern).Result()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("Found keys: %v\n", keys)
			stats, err := getPlayerStatsFromKeys(ctx, keys)
			fmt.Printf("Fetched stats: %v\n", stats)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			playerStats = append(playerStats, stats...)
		}
	} else {
		redisPattern := "player_stats:*"

		keys, err := rdb.Keys(ctx, redisPattern).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(keys) == 0 {
			allStats, err := fetchPlayerStatsFromJson()
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to fetch player stats from JSON: %v", err), http.StatusInternalServerError)
				return
			}
			keys, err = rdb.Keys(ctx, redisPattern).Result()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			playerStats = allStats
		} else {
			playerStats, err = getPlayerStatsFromKeys(ctx, keys)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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

func writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}
	return nil
}
