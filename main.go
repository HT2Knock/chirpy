package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
)

const (
	defaultDir  = "./"
	defaultPort = "8080"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	hitCount := "Hits: " + strconv.FormatInt(int64(cfg.fileServerHits.Load()), 10)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(hitCount))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)

	w.WriteHeader(200)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
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

	apiCfg := apiConfig{}

	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.metricHandler)
	mux.HandleFunc("POST /api/reset", apiCfg.resetHandler)
	mux.Handle("/app/", middlewareLog(apiCfg.middlewareMetricInc(http.StripPrefix("/app/", http.FileServer(http.Dir(absAssetDir))))))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", absAssetDir, port)
	log.Fatal(s.ListenAndServe())
}
