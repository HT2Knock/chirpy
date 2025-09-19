package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/T2Knock/chirpy/internal/auth"
	"github.com/T2Knock/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type requestCreateChirp struct {
	Body string `json:"body"`
}

type returnErr struct {
	Error string `json:"error"`
}

func filterProfanity(input string) string {
	profaneMap := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Fields(input)
	for i, word := range words {
		if _, found := profaneMap[strings.ToLower(word)]; found {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	request := requestCreateChirp{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding parameters: %s", err)

		writeJSON(w, http.StatusInternalServerError, returnErr{Error: "Something went wrong"})
		return
	}

	if len(request.Body) > 140 {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: "Chirp is too long"})
		return
	}

	userID := auth.UserIDFromContext(r.Context())

	user, err := cfg.dbQueries.GetUser(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: "User not found"})
		return
	}

	newChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: sql.NullString{String: filterProfanity(request.Body), Valid: true}, CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
		return
	}

	chirp := Chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body.String,
		UserID:    newChirp.UserID,
	}

	writeJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: fmt.Sprintf("%v", err)})
		return
	}

	chirps := make([]Chirp, 0, len(dbChirps))

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body.String,
			UserID:    dbChirp.UserID,
		})
	}

	writeJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: fmt.Sprintf("%v", err)})
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, returnErr{Error: fmt.Sprintf("%v", err)})
		return
	} else if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: fmt.Sprintf("%v", err)})
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body.String,
		UserID:    dbChirp.UserID,
	}

	writeJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: err.Error()})
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
		return
	}

	userID := auth.UserIDFromContext(r.Context())

	if chirp.UserID != userID {
		writeStatus(w, http.StatusForbidden)
		return
	}

	if err := cfg.dbQueries.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{ID: chirpID, UserID: userID}); err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
		return
	}

	writeStatus(w, http.StatusNoContent)
}
