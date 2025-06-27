package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	s := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
