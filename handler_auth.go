package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/T2Knock/chirpy/internal/auth"
	"github.com/T2Knock/chirpy/internal/database"
)

type requestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshToken struct {
	Token string `json:"token"`
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

	jwt, err := auth.MakeJWT(findUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
		return
	}

	if err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    findUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
		return
	}

	user := User{
		ID:           findUser.ID,
		CreatedAt:    findUser.CreatedAt,
		UpdatedAt:    findUser.UpdatedAt,
		Email:        findUser.Email,
		Token:        jwt,
		RefreshToken: refreshToken,
	}

	writeJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, returnErr{Error: err.Error()})
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, returnErr{Error: err.Error()})
			return
		}

		ctx := auth.WithUserID(r.Context(), userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: err.Error()})
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: err.Error()})
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		writeJSON(w, http.StatusUnauthorized, returnErr{Error: "Refresh token expired"})
		return
	}

	if refreshToken.RevokedAt.Valid {
		writeJSON(w, http.StatusUnauthorized, returnErr{Error: "Refresh token revoked"})
		return
	}

	jwt, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, RefreshToken{Token: jwt})
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: err.Error()})
		return
	}

	if err := cfg.dbQueries.UpdateRevokeRefreshToken(r.Context(), database.UpdateRevokeRefreshTokenParams{Token: token, RevokedAt: sql.NullTime{Time: time.Now(), Valid: true}}); err != nil {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: err.Error()})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	log.Printf("Revoked token %v \n", token)
}
