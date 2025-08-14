package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type requestCreateUser struct {
	Email string `json:"email"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	request := requestCreateUser{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding parameters: %s", err)

		writeJSON(w, 500, returnErr{Error: "Something went wrong"})
		return
	}
}
