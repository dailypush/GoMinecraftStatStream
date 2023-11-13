package main

import (
	"context"
	"fmt"
	"log"
	"strings"
)

// getAllStatTypes retrieves all unique statistic types from Redis.
// It returns a slice of stat types or an error.
func getAllStatTypes(ctx context.Context) ([]string, error) {
	playerNames, err := getAllPlayersFromRedis(ctx)
	if err != nil {
		log.Printf("Error retrieving players from Redis: %v", err)
		return nil, err
	}

	if len(playerNames) == 0 {
		return []string{}, nil // Return an empty slice if no players are found
	}

	statTypeSet := make(map[string]struct{})

	for _, playerName := range playerNames {
		pattern := fmt.Sprintf("player_stats:%s:*", playerName)
		var cursor uint64
		for {
			var keys []string
			keys, cursor, err = rdb.Scan(ctx, cursor, pattern, 50).Result() // Adjusted count for efficiency
			if err != nil {
				log.Printf("Error scanning Redis keys for player %s: %v", playerName, err)
				return nil, err
			}

			for _, key := range keys {
				statType := strings.TrimPrefix(key, fmt.Sprintf("player_stats:%s:", playerName))
				statTypeSet[statType] = struct{}{}
			}

			if cursor == 0 {
				break
			}
		}
	}

	allStatTypes := make([]string, 0, len(statTypeSet))
	for statType := range statTypeSet {
		allStatTypes = append(allStatTypes, statType)
	}

	return allStatTypes, nil
}
