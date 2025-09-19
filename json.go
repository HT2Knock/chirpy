package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if value != nil {
		if err := json.NewEncoder(w).Encode(value); err != nil {
			log.Printf("Error encoding JSON: %s", err)
		}
	}
}

func writeStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
