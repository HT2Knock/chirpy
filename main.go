package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	defaultDir  = "./"
	defaultPort = "8080"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	var assetDir, port string

	flag.StringVar(&assetDir, "asset-dir", defaultDir, "directory to serve files from")
	flag.StringVar(&port, "port", defaultPort, "port to serve on")
	flag.Parse()

	absAssetDir, err := filepath.Abs(assetDir)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path for asset directory '%s': %v", assetDir, err)
	}

	info, err := os.Stat(absAssetDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Asset directory does not exist: %s", absAssetDir)
		} else {
			log.Fatalf("Error checking asset directory '%s': %v", absAssetDir, err)
		}
	}
	if !info.IsDir() {
		log.Fatalf("Asset path '%s' is not a directory.", absAssetDir)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthHandler)
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(absAssetDir))))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", absAssetDir, port)
	log.Fatal(s.ListenAndServe())
}
