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
)

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

func parsePlayerList(playerList string) []string {
	trimmed := strings.TrimSpace(strings.TrimPrefix(playerList, "There are x/y players online:"))
	if trimmed == "" {
		return []string{}
	}

	return strings.Split(trimmed, ", ")
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

	playerNames := parsePlayerList(playerList)

	statTypes := []string{
		// List the stat types you want to fetch here, e.g.:
		"blocks_mined",
		"arrows_shot",
		// ...
	}

	stats := []PlayerStats{}

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

			err = storePlayerStatInRedis(stat)
			if err != nil {
				log.Printf("Error storing player stat in Redis: %v", err)
			} else {
				// log.Printf("Successfully set stat in Redis: Key=player_stats:%s:%s, Value=%d", playerName, statType, stat.Value)
			}

			stats = append(stats, stat)
		}
	}

	return stats
}

func processStatFile(path string) ([]PlayerStats, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read stat file: %v", err)
	}

	var rawStats struct {
		Stats map[string]map[string]int
	}
	err = json.Unmarshal(data, &rawStats)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stat file: %v", err)
	}

	fileName := filepath.Base(path)
	playerUUID := fileName[:len(fileName)-len(".json")]
	playerName, err := uuidToPlayerName(playerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert UUID to player name: %v", err)
	}

	stats := []PlayerStats{}

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
	log.Printf("PlayerStats from file %s: %+v\n", path, stats)
	return stats, nil
}

func fetchPlayerStatsFromJson() []PlayerStats {
	jsonFiles, err := ioutil.ReadDir(JsonStatsDirectory)
	if err != nil {
		log.Fatalf("failed to read JSON stats directory: %v", err)
	}

	ctx := context.Background()
	var allStats []PlayerStats
	for _, file := range jsonFiles {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(JsonStatsDirectory, file.Name())
		log.Printf("Updating stats from file: %s", filePath)
		stats, err := processStatFile(filePath)
		if err != nil {
			log.Printf("failed to read stats from file %q: %v", filePath, err)
			continue
		}

		for _, stat := range stats {
			key := fmt.Sprintf("player_stats:%s:%s", stat.Player, stat.StatType)
			err := rdb.Set(ctx, key, stat.Value, 0).Err()
			if err != nil {
				log.Printf("Failed to set stat in Redis: Key=%s, Value=%d", key, stat.Value)
			} else {
				log.Printf("Successfully set stat in Redis: Key=%s, Value=%d", key, stat.Value)
			}
		}
		allStats = append(allStats, stats...)
	}
	return allStats
}
