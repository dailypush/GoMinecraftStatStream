package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var ctx = context.Background()

func pollPlayerStats(interval time.Duration) {
	// Track the hashes of processed stat files
	hashes := make(map[string]string)

	for {
		// Check for updated stat files
		updatedFiles, err := getUpdatedStatFiles(hashes)
		if err != nil {
			log.Printf("Error checking for updated stat files: %v", err)
		}

		if len(updatedFiles) > 0 {
			// Fetch player stats from the updated files
			err = fetchPlayerStatsFromFiles(updatedFiles)
			if err != nil {
				log.Printf("Error fetching player stats from files: %v", err)
			}

			// Update the hashes of processed files
			for _, file := range updatedFiles {
				hashes[file] = getFileHash(file)
				log.Printf("Updated stats from file: %s", file)
			}
		}

		// Wait for the specified interval before checking for updates again
		time.Sleep(interval)
	}
}

func fetchPlayerStatsFromFiles(files []string) error {
	for _, file := range files {
		fileStats, err := processStatFile(file)
		if err != nil {
			return fmt.Errorf("failed to process stat file: %v", err)
		}

		for _, stat := range fileStats {
			err = storePlayerStatInRedis(stat)
			if err != nil {
				log.Printf("Error storing player stat in Redis: %v", err)
			} else {
				log.Printf("Successfully set stat in Redis: Key=player_stats:%s:%s, Value=%d", stat.Player, stat.StatType, stat.Value)
			}
		}
	}

	return nil
}
func getUpdatedStatFiles(hashes map[string]string) ([]string, error) {
	updatedFiles := []string{}

	err := filepath.Walk(JsonStatsDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// Get the file's hash and compare it with the previously recorded hash
			fileHash := getFileHash(path)
			if fileHash != hashes[path] {
				updatedFiles = append(updatedFiles, path)
			}
			hashes[path] = fileHash
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to check for updated stat files: %v", err)
	}

	return updatedFiles, nil
}

func getFileHash(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open file %s: %v", path, err)
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Printf("Failed to compute hash for file %s: %v", path, err)
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
