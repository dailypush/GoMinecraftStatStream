package main

import (
	"encoding/json"
	"fmt"
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
				log.Printf("Successfully set stat in Redis: Key=player_stats:%s:%s, Value=%d", playerName, statType, stat.Value)
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

	return stats, nil
}

func fetchPlayerStatsFromJson() []PlayerStats {
	stats := []PlayerStats{}

	err := filepath.Walk(JsonStatsDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileStats, err := processStatFile(path)
			if err != nil {
				log.Printf("Error processing stat file %s: %v", path, err)
				return nil
			}

			for _, stat := range fileStats {
				err = storePlayerStatInRedis(stat)
				if err != nil {
					log.Printf("Error storing player stat in Redis: %v", err)
				} else {
					log.Printf("Successfully set stat in Redis: Key=player_stats:%s:%s, Value=%d", stat.Player, stat.StatType, stat.Value)
				}
			}

			stats = append(stats, fileStats...)
		}
		return nil
	})

	if err != nil {
		log.Fatal("Failed to fetch stats from JSON files:", err)
	}

	return stats
}
