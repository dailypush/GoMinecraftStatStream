package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStats(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getStats)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var stats []PlayerStats
	err = json.Unmarshal(rr.Body.Bytes(), &stats)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, stats, "Expected non-empty stats")
}
