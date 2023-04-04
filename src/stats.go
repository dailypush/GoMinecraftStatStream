package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	// List the stat types you want to fetch here
	statTypes := []string{
		// ...
	}

	err := filepath.Walk(JsonStatsDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			playerStats, err := processStatFile(path, info, statTypes)
			if err != nil {
				log.Printf("Failed to process stat file: %v", err)
			} else {
				stats = append(stats, playerStats...)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("Failed to fetch stats from JSON files:", err)
	}

	return stats
}

func processStatFile(path string, info os.FileInfo, statTypes []string) ([]PlayerStats, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read stat file: %w", err)
	}

	var rawStats map[string]map[string]int
	err = json.Unmarshal(data, &rawStats)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stat file: %w", err)
	}

	playerUUID := info.Name()[:len(info.Name())-len(".json")]
	playerName, err := uuidToPlayerName(playerUUID)
	if err != nil {
		log.Printf("Failed to convert UUID to player name: %v", err)
		return nil, err
	}

	stats := []PlayerStats{}
	for _, statType := range statTypes {
		if category, ok := rawStats[statType]; ok {
			for stat, value := range category {
				stat := PlayerStats{
					Player:   playerName,
					StatType: stat,
					Value:    value,
				}

				stats = append(stats, stat)
			}
		}
	}

	return stats, nil
}
