package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.Fields(authHeader)

	if len(parts) != 2 || strings.ToLower(parts[0]) != "apikey" {
		return "", fmt.Errorf("invalid authorization header ")
	}

	apiKey := parts[1]
	if apiKey == "" {
		return "", fmt.Errorf("api key not found in authorization header")
	}

	return apiKey, nil
}
