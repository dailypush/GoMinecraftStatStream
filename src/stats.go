package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gorcon/rcon"
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
	switch StatsSource {
	case "rcon":
		return fetchPlayerStatsFromRcon()
	case "json":
		return fetchPlayerStatsFromJson()
	default:
		log.Fatalf("Invalid stats source: %s", StatsSource)
		return nil
	}
}

func fetchPlayerStatsFromRcon() []PlayerStats {
	conn, err := rcon.Dial(RconAddress, RconPassword)
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

func fetchPlayerStatsFromJson() []PlayerStats {
	stats := []PlayerStats{}

	err := filepath.Walk(JsonStatsDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Failed to read stat file: %v", err)
				return nil
			}

			var rawStats struct {
				Stats map[string]map[string]int
			}
			err = json.Unmarshal(data, &rawStats)
			if err != nil {
				log.Printf("Failed to parse stat file: %v", err)
				return nil
			}

			playerUUID := info.Name()[:len(info.Name())-len(".json")]
			playerName, err := uuidToPlayerName(playerUUID)
			if err != nil {
				log.Printf("Failed to convert UUID to player name: %v", err)
				return nil
			}

			// Iterate over all stat types in the JSON file
			for statType, category := range rawStats.Stats {
				for stat, value := range category {
					playerStat := PlayerStats{
						Player:   playerName,
						StatType: fmt.Sprintf("%s:%s", statType, stat),
						Value:    value,
					}

					stats = append(stats, playerStat)
				}
			}

		}
		return nil
	})

	if err != nil {
		log.Fatal("Failed to fetch stats from JSON files:", err)
	}

	return stats
}
