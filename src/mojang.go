package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MojangProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func uuidToPlayerName(uuid string) (string, error) {
	// Remove dashes from UUID
	strippedUUID := strings.ReplaceAll(uuid, "-", "")

	// Make an API request to Mojang to fetch the player name
	resp, err := http.Get(fmt.Sprintf("https://api.mojang.com/user/profiles/%s/names", strippedUUID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get player name, status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var profile MojangProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return "", err
	}

	return profile.Name, nil
}
