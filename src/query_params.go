package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type QueryParams struct {
	PlayerNames []string
	StatType    string
	GroupBy     string
	SortOrder   string
	Top         int
	Category    string
}

func parseQueryParams(r *http.Request) (QueryParams, error) {
	query := r.URL.Query()

	playerNames := strings.Split(query.Get("playerNames"), ",")
	if len(playerNames) == 1 && playerNames[0] == "" {
		playerNames = []string{}
	}

	top, err := strconv.Atoi(query.Get("top"))
	if err != nil && query.Get("top") != "" {
		return QueryParams{}, fmt.Errorf("top parameter must be a valid integer")
	}

	return QueryParams{
		PlayerNames: playerNames,
		Category:    query.Get("category"),
		GroupBy:     query.Get("groupBy"),
		SortOrder:   query.Get("sortOrder"),
		Top:         top,
	}, nil
}
