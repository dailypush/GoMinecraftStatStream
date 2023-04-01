package main

import (
	"context"
	"fmt"

	"github.com/gorcon/rcon"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func fetchPlayerStats() []PlayerStats {
	conn, err := rcon.Dial("your.minecraftserver.com:25575", "your_rcon_password")
	if err != nil {
		log.Fatal("Could not connect to the Minecraft server:", err)
	}
	defer conn.Close()

	stats := []PlayerStats{
		// Add sample player stats
		{
			Player:   "Steve",
			StatType: "blocks_mined",
			Value:    100,
		},
		{
			Player:   "Alex",
			StatType: "arrows_shot",
			Value:    50,
		},
	}

	ctx := context.Background()

	for _, stat := range stats {
		key := fmt.Sprintf("player_stats:%s:%s", stat.Player, stat.StatType)

		err := rdb.Set(ctx, key, stat.Value, 0).Err()
		if err != nil {
			log.Printf("Failed to set stat in Redis: %v", err)
		}
	}

	return stats
}
