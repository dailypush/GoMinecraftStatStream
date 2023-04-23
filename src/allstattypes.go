package main

import (
	"fmt"
	"strings"
)

func getAllStatTypes() ([]string, error) {
	playerNames, err := getAllPlayersFromRedis()
	if err != nil {
		return nil, err
	}

	statTypeSet := make(map[string]struct{})

	for _, playerName := range playerNames {
		pattern := fmt.Sprintf("player_stats:%s:*", playerName)

		var cursor uint64
		for {
			var keys []string
			keys, cursor, err = rdb.Scan(ctx, cursor, pattern, 10).Result()
			if err != nil {
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
