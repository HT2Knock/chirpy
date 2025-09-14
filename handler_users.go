package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/T2Knock/chirpy/internal/auth"
	"github.com/T2Knock/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type requestCreateUser struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	request := requestCreateUser{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding parameters: %s", err)

		writeJSON(w, http.StatusInternalServerError, returnErr{Error: "Something went wrong"})
		return
	}

	if request.Email == "" || request.Password == "" {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: "Missing required parameters"})
		return
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: fmt.Sprintf("%v", err)})
	}

	createdUser, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: request.Email, HashedPassword: hashedPassword, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: fmt.Sprintf("%v", err)})
	}

	user := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	writeJSON(w, http.StatusCreated, user)
}
