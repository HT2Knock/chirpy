package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/T2Knock/chirpy/internal/database"
	"github.com/google/uuid"
)

type requestPolkaWebhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) webhookHandler(w http.ResponseWriter, r *http.Request) {
	request := requestPolkaWebhook{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding parameters: %s", err)

		writeJSON(w, http.StatusInternalServerError, returnErr{Error: "Something went wrong"})
		return
	}

	if request.Event != "user.upgraded" {
		writeStatus(w, http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(request.Data.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
		return
	}

	if err := cfg.dbQueries.UpdateChirpyRed(r.Context(), database.UpdateChirpyRedParams{ID: userID, IsChirpyRed: true}); err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
		return
	}

	writeStatus(w, http.StatusNoContent)
}
