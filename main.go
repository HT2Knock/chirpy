package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	defaultDir  = "./assets"
	defaultPort = "8080"
)

// getConfigValue checks flag, then environment variable, then default value.
func getConfigValue(flagVal, envVarName, defaultVal string) string {
	if flagVal != "" {
		return flagVal
	}
	if envVal := os.Getenv(envVarName); envVal != "" {
		return envVal
	}
	return defaultVal
}

func main() {
	var assetDir, port string

	flag.StringVar(&assetDir, "asset-dir", "", "directory to serve files from")
	flag.StringVar(&port, "port", "", "port to serve on")
	flag.Parse()

	assetDir = getConfigValue(assetDir, "ASSET_DIR", defaultDir)
	port = getConfigValue(port, "PORT", defaultPort)

	// Resolve absolute path for clarity and consistency
	absAssetDir, err := filepath.Abs(assetDir)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path for asset directory '%s': %v", assetDir, err)
	}

	// Validate asset directory exists and is a directory
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
	mux.Handle("/", http.FileServer(http.Dir(absAssetDir)))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", absAssetDir, port)
	log.Fatal(s.ListenAndServe())
}
