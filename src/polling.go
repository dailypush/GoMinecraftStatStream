package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

var ctx = context.Background()

func pollPlayerStats(interval time.Duration) {
	for {
		stats := fetchPlayerStats()

		// Update the stats in Redis
		for _, stat := range stats {
			key := fmt.Sprintf("player_stats:%s:%s", stat.Player, stat.StatType)
			err := rdb.Set(ctx, key, stat.Value, 0).Err()
			if err != nil {
				log.Printf("Failed to update stat in Redis: %v", err)
			}
		}

		// Wait for the specified interval before fetching the stats again
		time.Sleep(interval)
	}
}
