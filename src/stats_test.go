package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchPlayerStats(t *testing.T) {
	stats := fetchPlayerStats()

	assert.NotEmpty(t, stats, "Expected non-empty stats")

	for _, stat := range stats {
		assert.NotEmpty(t, stat.Player, "Expected player name to be non-empty")
		assert.NotEmpty(t, stat.StatType, "Expected stat type to be non-empty")
		assert.True(t, stat.Value >= 0, "Expected stat value to be non-negative")
	}
}
