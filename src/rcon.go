package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gorcon/rcon"
)

func getPlayerStat(conn *rcon.Conn, player, statType string) (int, error) {
	// Replace the following line with the appropriate command to fetch the player's stat
	response, err := conn.Execute(fmt.Sprintf("stats %s %s", player, statType))

	if err != nil {
		return 0, err
	}

	// Parse the response to extract the stat value (replace this with the correct parsing logic)
	valueStr := strings.TrimSpace(response) // Adjust this if the response format requires more complex parsing

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse stat value: %v", err)
	}

	return value, nil
}
