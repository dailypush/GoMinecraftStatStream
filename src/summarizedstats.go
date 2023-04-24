package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func getSummarizedStats(statTypes []string) (map[string]map[string]int64, map[string]map[string]int64, error) {
	playerNames, err := getAllPlayersFromRedis()
	if err != nil {
		return nil, nil, err
	}

	aggregatedStats := make(map[string]map[string]int64)
	individualStats := make(map[string]map[string]int64)

	for _, playerName := range playerNames {
		for _, statType := range statTypes {
			pattern := fmt.Sprintf("player_stats:%s:*%s*", playerName, statType)

			var cursor uint64
			for {
				var keys []string
				keys, cursor, err = rdb.Scan(ctx, cursor, pattern, 10).Result()
				if err != nil {
					return nil, nil, err
				}

				for _, key := range keys {
					stat, err := rdb.Get(ctx, key).Int64()

					if err == redis.Nil {
						stat = 0
					} else if err != nil {
						return nil, nil, err
					}

					if individualStats[playerName] == nil {
						individualStats[playerName] = make(map[string]int64)
					}
					individualStats[playerName][key] = stat

					if aggregatedStats[playerName] == nil {
						aggregatedStats[playerName] = make(map[string]int64)
					}
					aggregatedStats[playerName][statType] += stat
				}

				if cursor == 0 {
					break
				}
			}
		}
	}

	return aggregatedStats, individualStats, nil
}
