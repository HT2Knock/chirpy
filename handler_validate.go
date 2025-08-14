package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type returnErr struct {
	Error string `json:"error"`
}

type returnVals struct {
	CleanBody string `json:"cleaned_body"`
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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	params := parameters{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)

		writeJSON(w, 500, returnErr{Error: "Something went wrong"})
		return
	}

	if len(params.Body) > 140 {
		writeJSON(w, 400, returnErr{Error: "Chirp is too long"})
		return
	}

	writeJSON(w, 200, returnVals{CleanBody: filterProfanity(params.Body)})
}
