package main

import (
	"os"
	"time"
)

var (
	StatsSource        string
	RconAddress        string
	RconPassword       string
	JsonStatsDirectory string
	ServerPort         string
	PollingInterval    time.Duration
)

func init() {
	StatsSource = getEnv("STATS_SOURCE", "rcon")
	RconAddress = getEnv("RCON_ADDRESS", "localhost:25575")
	RconPassword = getEnv("RCON_PASSWORD", "your_rcon_password")
	JsonStatsDirectory = getEnv("JSON_STATS_DIRECTORY", "./json_stats")
	PollingInterval = getDurationEnv("POLLING_INTERVAL", "5m")
	ServerPort = getEnv("SERVER_PORT", "8080")
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getDurationEnv(key, fallback string) time.Duration {
	value := os.Getenv(key)
	if len(value) == 0 {
		value = fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}
	return duration
}
