package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("invalid authorization header")
	}
	apiKey := strings.TrimSpace(strings.TrimPrefix(authHeader, "ApiKey "))
	if apiKey == "" {
		return "", fmt.Errorf("empty API key")
	}
	return apiKey, nil
}
