package main

import (
	"encoding/json"
	"fmt"
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

	jwt, err := auth.MakeJWT(findUser.ID, cfg.jwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
	}

	sixtyDays := time.Hour * 24 * 60
	if err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refreshToken, UserID: findUser.ID, ExpiresAt: time.Now().Add(sixtyDays), CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != nil {
		writeJSON(w, http.StatusInternalServerError, returnErr{Error: err.Error()})
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
		fmt.Println(userID)
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

	jwt, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, returnErr{Error: err.Error()})
	}

	writeJSON(w, http.StatusOK, RefreshToken{Token: jwt})
}
