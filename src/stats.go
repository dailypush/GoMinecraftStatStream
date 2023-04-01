package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gorcon/rcon"
	"github.com/go-redis/redis/v8"
)

type PlayerStats struct {
	Player   string `json:"player"`
	StatType string `json:"statType"`
	Value    int    `json:"value"`
}

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func getPlayerStat(conn *rcon.Conn, player, statType string) (int, error) {
	// Replace the following line with the appropriate command to fetch the player's stat
	response, err := conn.Execute(fmt.Sprintf("stats %s %s", player, statType))

	if err != nil {
		return 0, err
	}

	// Parse the response to extract the stat value (replace this with the correct parsing logic)
	value := 0 // Set this to the actual value extracted from the response

	return value, nil
}

func fetchPlayerStats() []PlayerStats {
	conn, err := rcon.Dial("your.minecraftserver.com:25575", "your_rcon_password")
	if err != nil {
		log.Fatal("Could not connect to the Minecraft server:", err)
	}
	defer conn.Close()

	playerList, err := conn.Execute("list")
	if err != nil {
		log.Fatal("Could not fetch player list:", err)
	}

	playerNames := strings.Split(strings.TrimSpace(strings.TrimPrefix(playerList, "There are x/y players online:")), ", ")

	statTypes := []string{
		// List the stat types you want to fetch here, e.g.:
		"blocks_mined",
		"arrows_shot",
		// ...
	}

	stats := []PlayerStats{}

	ctx := context.Background()

	for _, playerName := range playerNames {
		for _, statType := range statTypes {
			value, err := getPlayerStat(conn, playerName, statType)
			if err != nil {
				log.Printf("Failed to get stat '%s' for player '%s': %v", statType, playerName, err)
				continue
			}

			stat := PlayerStats{
				Player:   playerName,
				StatType: statType,
				Value:    value,
			}

			key := fmt.Sprintf("player_stats:%s:%s", stat.Player, stat.StatType)
			err = rdb.Set(ctx, key, stat.Value, 0).Err()
			if err != nil {
				log.Printf("Failed to set stat in Redis: %v", err)
			}

			stats = append(stats, stat)
		}
	}

	return stats
}
