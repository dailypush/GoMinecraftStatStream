package main

import (
	"context"
	"fmt"
	"log"
	"strings"
)

func getPlayerStatsFromKeys(ctx context.Context, keys []string) ([]PlayerStats, error) {
	var playerStats []PlayerStats

	log.Printf("Found %d keys in Redis.", len(keys))
	for _, key := range keys {
		value, err := rdb.Get(ctx, key).Int64()
		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %q: %v", key, err)
		}
		keyParts := strings.SplitN(key, ":", 3)
		playerName := keyParts[1]
		statType := keyParts[2]
		playerStat := PlayerStats{
			Player:   playerName,
			StatType: statType,
			Value:    int(value),
		}
		playerStats = append(playerStats, playerStat)
		log.Printf("Retrieved player stat from Redis: Key=%s, Value=%d", key, value)
	}
	return playerStats, nil
}
