package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

type MojangProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getPlayerNameFromUUID(uuid string) (string, error) {
	// Try to get the username from Redis
	key := fmt.Sprintf("uuid_to_username:%s", uuid)
	username, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		// If the username is not in Redis, query Mojang's API
		username, err = getMojangUsername(uuid)
		if err != nil {
			return "", fmt.Errorf("failed to get username from Mojang API: %v", err)
		}

		// Store the UUID to username mapping in Redis
		err = storeUuidToUsernameMappingInRedis(uuid, username)
		if err != nil {
			log.Printf("Error storing uuid to username mapping in Redis: %v", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("failed to get username from Redis: %v", err)
	}

	return username, nil
}

func getMojangUsername(uuid string) (string, error) {
	// Remove dashes from UUID
	strippedUUID := strings.ReplaceAll(uuid, "-", "")

	// Make an API request to Mojang to fetch the player name
	resp, err := http.Get(fmt.Sprintf("https://api.mojang.com/user/profile/%s", strippedUUID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get player name, status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var profile MojangProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return "", err
	}

	return profile.Name, nil
}
