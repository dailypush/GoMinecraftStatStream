package main

import (
	"os"
)

var (
	StatsSource        string
	RconAddress        string
	RconPassword       string
	JsonStatsDirectory string
)

func init() {
	StatsSource = getEnv("STATS_SOURCE", "rcon")
	RconAddress = getEnv("RCON_ADDRESS", "localhost:25575")
	RconPassword = getEnv("RCON_PASSWORD", "your_rcon_password")
	JsonStatsDirectory = getEnv("JSON_STATS_DIRECTORY", "./json_stats")
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
