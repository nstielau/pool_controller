// Package api provides HTTP handlers for the pool controller.
package api

import (
	"net/http"
	"os"
	"regexp"
	"strings"
)

// AuthMiddleware validates Bearer tokens against TOKEN_REGEX environment variable.
func AuthMiddleware(next http.Handler) http.Handler {
	tokenRegex := os.Getenv("TOKEN_REGEX")
	if tokenRegex == "" {
		tokenRegex = ".*" // Default: accept any token
	}

	pattern, err := regexp.Compile(tokenRegex)
	if err != nil {
		// If regex is invalid, use permissive pattern
		pattern = regexp.MustCompile(".*")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			authHeader = r.Header.Get("Authentication") // Legacy support
		}

		token := ""
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 {
				token = parts[1]
			}
		}

		// Validate token against regex
		if !pattern.MatchString(token) {
			w.WriteHeader(http.StatusOK) // Original Python returned 200 with "Unauthed"
			w.Write([]byte("Unauthed"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
