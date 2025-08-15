package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/T2Knock/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	defaultDir  = "./"
	defaultPort = "8080"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
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

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file ")
	}

	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error open connection to database: %v", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       platform,
	}

	mux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricInc(http.StripPrefix("/app/", http.FileServer(http.Dir(absAssetDir))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", healthHandler)

	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpHandler)

	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", absAssetDir, port)
	log.Fatal(server.ListenAndServe())
}
