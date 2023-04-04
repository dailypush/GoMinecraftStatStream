package main

import (
	"fmt"

	"github.com/gorcon/rcon"
)

func getPlayerStat(conn *rcon.Conn, player, statType string) (int, error) {
	// Replace the following line with the appropriate command to fetch the player's stat
	response, err := conn.Execute(fmt.Sprintf("stats %s %s", player, statType))

	if err != nil {
		return 0, err
	}

	// Parse the response to extract the stat value (replace this with the correct parsing logic)
	value := 0 // Set this to the actual value extracted from the response

	return value, nil
}
