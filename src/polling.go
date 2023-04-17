package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

var ctx = context.Background()

func pollPlayerStats(interval time.Duration) {
	lastModified := time.Time{}

	for {
		// Check the modification time of the file
		fileInfo, err := os.Stat(JsonStatsDirectory)
		if err != nil {
			log.Printf("Error getting stats directory info: %v", err)
		} else if !fileInfo.ModTime().After(lastModified) {
			// Skip this iteration if the file hasn't been modified since the last iteration
			time.Sleep(interval)
			continue
		} else {
			lastModified = fileInfo.ModTime()
		}

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
