package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/T2Knock/chirpy/internal/auth"
)

type requestLogin struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	request := requestLogin{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding parameters: %s", err)

		writeJSON(w, http.StatusInternalServerError, returnErr{Error: "Something went wrong"})
		return
	}

	findUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, returnErr{Error: "Incorrect email or password"})
		return
	}

	if err := auth.CheckPasswordHash(request.Password, findUser.HashedPassword); err != nil {
		writeJSON(w, http.StatusUnauthorized, returnErr{Error: "Incorrect email or password"})
		return
	}

	jwt, err := auth.MakeJWT(findUser.ID, cfg.jwtSecret, time.Duration(request.ExpiresInSeconds)*time.Second)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: fmt.Sprintf("%s", err)})
	}

	user := User{
		ID:        findUser.ID,
		CreatedAt: findUser.CreatedAt,
		UpdatedAt: findUser.UpdatedAt,
		Email:     findUser.Email,
		Token:     jwt,
	}

	writeJSON(w, http.StatusOK, user)
}
