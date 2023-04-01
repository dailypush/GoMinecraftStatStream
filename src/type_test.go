package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerStats(t *testing.T) {
	playerStat := PlayerStats{
		Player:   "Steve",
		StatType: "blocks_mined",
		Value:    100,
	}

	assert.Equal(t, "Steve", playerStat.Player, "Expected player name to be 'Steve'")
	assert.Equal(t, "blocks_mined", playerStat.StatType, "Expected stat type to be 'blocks_mined'")
	assert.Equal(t, 100, playerStat.Value, "Expected stat value to be 100")
}
