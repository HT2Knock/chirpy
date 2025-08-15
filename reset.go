package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		writeJSON(w, http.StatusForbidden, returnErr{Error: "Forbidden"})
		return
	}

	cfg.fileServerHits.Store(0)
	if err := cfg.dbQueries.DeleteUsers(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: fmt.Sprintf("%v", err)})
		return
	}

	log.Println("Delete all users in database")
	w.WriteHeader(http.StatusOK)
}
