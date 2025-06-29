package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	defaultDir := "./assets"
	defaultPort := "8080"
	var assetDir, port string

	flag.StringVar(&assetDir, "asset-dir", "", "directory to serve files from")
	flag.StringVar(&port, "port", "", "port to serve on")
	flag.Parse()

	// Check flags, then env, then defaults
	if assetDir == "" {
		assetDir = os.Getenv("ASSET_DIR")
		if assetDir == "" {
			assetDir = defaultDir
		}
	}
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = defaultPort
		}
	}

	// Validate asset directory exists
	if info, err := os.Stat(assetDir); err != nil || !info.IsDir() {
		log.Fatalf("Invalid asset directory: %s", assetDir)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(assetDir)))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", assetDir, port)
	log.Fatal(s.ListenAndServe())
}
