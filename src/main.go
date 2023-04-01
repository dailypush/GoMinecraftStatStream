package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/stats", getStats)
	router.HandleFunc("/ws", wsEndpoint)

	log.Fatal(http.ListenAndServe(":8080", router))
}
